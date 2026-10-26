package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const cookieName = "mirmilowo_session"
const sessionTTL = 8 * time.Hour

type App struct {
	DB            *sql.DB
	Origin        string
	SecureCookies bool
	dummyHash     []byte
}

type User struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func New(db *sql.DB, origin string, secure bool) *App {
	// Equal-cost password verification for unknown accounts.
	hash, err := bcrypt.GenerateFromPassword([]byte("not-a-user-password"), 12)
	if err != nil {
		panic(err)
	}
	return &App{DB: db, Origin: origin, SecureCookies: secure, dummyHash: hash}
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", a.health)
	mux.HandleFunc("GET /api/auth/me", a.me)
	mux.HandleFunc("POST /api/auth/login", a.mutation(a.login))
	mux.HandleFunc("POST /api/auth/logout", a.mutation(a.logout))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		mux.ServeHTTP(w, r)
	})
}

// JSON plus exact Origin validation also protects login/logout against CSRF.
func (a *App) mutation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != a.Origin {
			fail(w, http.StatusForbidden, "Niedozwolone źródło żądania.")
			return
		}
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			fail(w, http.StatusUnsupportedMediaType, "Wymagany format JSON.")
			return
		}
		next(w, r)
	}
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := a.DB.PingContext(ctx); err != nil {
		fail(w, http.StatusServiceUnavailable, "Baza danych jest niedostępna.")
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		fail(w, http.StatusBadRequest, "Nieprawidłowe dane formularza.")
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		fail(w, http.StatusBadRequest, "Nieprawidłowe dane formularza.")
		return
	}
	input.Email = strings.TrimSpace(input.Email)
	// Passwords are never trimmed; spaces can be part of the password.
	if input.Email == "" || len(input.Email) > 254 || input.Password == "" || len(input.Password) > 72 {
		fail(w, http.StatusBadRequest, "Podaj adres e-mail i hasło (maksymalnie 72 bajty).")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user User
	var hash string
	err := a.DB.QueryRowContext(ctx,
		"SELECT id, email, password_hash, role FROM users WHERE email = ? LIMIT 1", input.Email,
	).Scan(&user.ID, &user.Email, &hash, &user.Role)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.dbError(w, err)
		return
	}
	passwordHash := []byte(hash)
	if errors.Is(err, sql.ErrNoRows) {
		passwordHash = a.dummyHash
	}
	verified := bcrypt.CompareHashAndPassword(passwordHash, []byte(input.Password)) == nil
	if err != nil || !verified || (user.Role != "admin" && user.Role != "mieszkaniec") {
		fail(w, http.StatusUnauthorized, "Nieprawidłowy adres e-mail lub hasło.")
		return
	}
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		a.dbError(w, err)
		return
	}
	token := hex.EncodeToString(bytes[:])
	expires := time.Now().UTC().Add(sessionTTL)
	// Rotate the current browser session atomically after successful authentication.
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		a.dbError(w, err)
		return
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)", tokenHash(token), user.ID, expires); err != nil {
		a.dbError(w, err)
		return
	}
	if old, err := r.Cookie(cookieName); err == nil {
		if _, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash(old.Value)); err != nil {
			a.dbError(w, err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		a.dbError(w, err)
		return
	}
	http.SetCookie(w, a.cookie(token, expires, int(sessionTTL.Seconds())))
	respond(w, http.StatusOK, map[string]any{"user": user})
}

func (a *App) me(w http.ResponseWriter, r *http.Request) {
	user, err := a.currentUser(r)
	if err != nil {
		a.dbError(w, err)
		return
	}
	if user == nil {
		http.SetCookie(w, a.cookie("", time.Unix(1, 0), -1))
	}
	respond(w, http.StatusOK, map[string]any{"user": user})
}

func (a *App) currentUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || len(cookie.Value) != 64 {
		return nil, nil
	}
	if _, err := hex.DecodeString(cookie.Value); err != nil {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user User
	err = a.DB.QueryRowContext(ctx, `SELECT u.id, u.email, u.role FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = ? AND s.expires_at > UTC_TIMESTAMP()`, tokenHash(cookie.Value)).Scan(&user.ID, &user.Email, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if _, err = a.DB.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash(cookie.Value)); err != nil {
			a.dbError(w, err)
			return
		}
	}
	http.SetCookie(w, a.cookie("", time.Unix(1, 0), -1))
	respond(w, http.StatusOK, map[string]bool{"ok": true})
}

// RequireRole is ready for future private API handlers; UI state is not authorization.
func (a *App) RequireRole(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := a.currentUser(r)
		if err != nil {
			a.dbError(w, err)
			return
		}
		if user == nil {
			fail(w, http.StatusUnauthorized, "Zaloguj się, aby kontynuować.")
			return
		}
		for _, role := range roles {
			if user.Role == role {
				next(w, r)
				return
			}
		}
		fail(w, http.StatusForbidden, "Brak uprawnień.")
	}
}

func (a *App) cookie(value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{Name: cookieName, Value: value, Path: "/", HttpOnly: true, Secure: a.SecureCookies, SameSite: http.SameSiteLaxMode, Expires: expires, MaxAge: maxAge}
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
func (a *App) dbError(w http.ResponseWriter, err error) {
	log.Printf("API error: %v", err)
	fail(w, http.StatusServiceUnavailable, "Usługa jest chwilowo niedostępna. Spróbuj ponownie później.")
}
func fail(w http.ResponseWriter, status int, message string) {
	respond(w, status, map[string]string{"error": message})
}
func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
