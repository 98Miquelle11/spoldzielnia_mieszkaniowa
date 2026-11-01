package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const cookieName = "mirmilowo_session"
const sessionSeconds = 8 * 60 * 60

type User struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func tokenHash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
func sessionToken(r *http.Request) string {
	c, err := r.Cookie(cookieName)
	if err != nil || len(c.Value) != 64 {
		return ""
	}
	if _, err = hex.DecodeString(c.Value); err != nil {
		return ""
	}
	return c.Value
}
func (a *App) setCookie(w http.ResponseWriter, token string, age int) {
	expires := time.Now().Add(time.Duration(age) * time.Second)
	if age < 0 {
		expires = time.Unix(1, 0)
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: a.SecureCookies, MaxAge: age, Expires: expires})
}
func (a *App) currentUser(w http.ResponseWriter, r *http.Request) (*User, error) {
	token := sessionToken(r)
	if token == "" {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user User
	err := a.DB.QueryRowContext(ctx, `SELECT u.id, u.email, u.role FROM sessions s JOIN users u ON u.id = s.user_id WHERE s.token_hash = ? AND s.expires_at > UTC_TIMESTAMP(6)`, tokenHash(token)).Scan(&user.ID, &user.Email, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		a.setCookie(w, "", -1)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (a *App) me(w http.ResponseWriter, r *http.Request) {
	user, err := a.currentUser(w, r)
	if err != nil {
		a.dbError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"user": user})
}
func (a *App) RequireRole(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := a.currentUser(w, r)
		if err != nil {
			a.dbError(w, err)
			return
		}
		if user == nil {
			fail(w, 401, "Zaloguj się ponownie.")
			return
		}
		for _, role := range roles {
			if user.Role == role {
				next(w, r)
				return
			}
		}
		fail(w, 403, "Brak uprawnień administratora.")
	}
}
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := decodeStrictJSON(r.Body, &input); err != nil {
		fail(w, 400, "Nieprawidłowe dane logowania.")
		return
	}
	input.Email = strings.TrimSpace(input.Email)
	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email || len(input.Email) > 254 || input.Password == "" || len(input.Password) > 72 {
		fail(w, 400, "Podaj poprawny e-mail i hasło.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	var user User
	var hash string
	err = a.DB.QueryRowContext(ctx, "SELECT id, email, password_hash, role FROM users WHERE email = ? LIMIT 1", input.Email).Scan(&user.ID, &user.Email, &hash, &user.Role)
	missing := errors.Is(err, sql.ErrNoRows)
	if err != nil && !missing {
		a.dbError(w, err)
		return
	}
	if missing {
		hash = string(a.dummyHash)
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password))
	if missing || passwordErr != nil {
		fail(w, 401, "Nieprawidłowy e-mail lub hasło.")
		return
	}
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil {
		fail(w, 500, "Nie udało się utworzyć sesji.")
		return
	}
	token := hex.EncodeToString(bytes)
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		a.dbError(w, err)
		return
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)", tokenHash(token), user.ID, time.Now().UTC().Add(sessionSeconds*time.Second))
	if err != nil {
		a.dbError(w, err)
		return
	}
	if old := sessionToken(r); old != "" {
		if _, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash(old)); err != nil {
			a.dbError(w, err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		a.dbError(w, err)
		return
	}
	a.setCookie(w, token, sessionSeconds)
	writeJSON(w, 200, map[string]any{"user": user})
}
func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	var body struct{}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := decodeStrictJSON(r.Body, &body); err != nil {
		fail(w, 400, "Wymagany pusty obiekt JSON.")
		return
	}
	if token := sessionToken(r); token != "" {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if _, err := a.DB.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash(token)); err != nil {
			a.dbError(w, err)
			return
		}
	}
	a.setCookie(w, "", -1)
	writeJSON(w, 200, map[string]bool{"ok": true})
}
