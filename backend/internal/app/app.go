package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"mirmilowo/backend/internal/events"
)

type App struct {
	DB            *sql.DB
	Origin        string
	SecureCookies bool
	events        *events.Repository
	dummyHash     []byte
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
func fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
func (a *App) dbError(w http.ResponseWriter, err error) {
	log.Printf("database error: %v", err)
	fail(w, http.StatusServiceUnavailable, "Baza danych jest chwilowo niedostępna.")
}
func New(db *sql.DB, origin string, secure bool) *App {
	hash, err := bcrypt.GenerateFromPassword([]byte("not-a-user-password"), 12)
	if err != nil {
		panic(err)
	}
	return &App{DB: db, Origin: origin, SecureCookies: secure, events: events.NewRepository(db), dummyHash: hash}
}
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", a.health)
	mux.HandleFunc("GET /api/auth/me", a.me)
	mux.HandleFunc("POST /api/auth/login", a.mutation(a.login))
	mux.HandleFunc("POST /api/auth/logout", a.mutation(a.logout))
	mux.HandleFunc("GET /api/events", a.listPublishedEvents)
	mux.HandleFunc("GET /api/events/{id}", a.getPublishedEvent)
	mux.HandleFunc("GET /api/admin/events", a.RequireRole([]string{"admin"}, a.listAdminEvents))
	mux.HandleFunc("GET /api/admin/events/{id}", a.RequireRole([]string{"admin"}, a.getAdminEvent))
	mux.HandleFunc("POST /api/admin/events", a.mutation(a.RequireRole([]string{"admin"}, a.createAdminEvent)))
	mux.HandleFunc("PUT /api/admin/events/{id}", a.mutation(a.RequireRole([]string{"admin"}, a.updateAdminEvent)))
	mux.HandleFunc("DELETE /api/admin/events/{id}", a.mutation(a.RequireRole([]string{"admin"}, a.deleteAdminEvent)))
	return mux
}
func (a *App) mutation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != a.Origin || a.Origin == "" {
			fail(w, 403, "Niedozwolone źródło żądania.")
			return
		}
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			fail(w, 415, "Wymagany Content-Type: application/json.")
			return
		}
		next(w, r)
	}
}
func parseEventID(raw string) (int64, error) {
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid id")
		}
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
func decodeStrictJSON(r io.Reader, dst any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return errors.New("expected one JSON value")
	}
	return nil
}
func (a *App) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]bool{"ok": true})
}
