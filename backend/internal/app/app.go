package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"mirmilowo/backend/internal/events"
)

type App struct {
	DB            *sql.DB
	Origin        string
	SecureCookies bool

	events    *events.Repository
	dummyHash []byte
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func New(db *sql.DB, origin string, secure bool) *App {
	// Używany przy logowaniu na nieistniejące konto,
	// aby również wtedy wykonać porównanie bcrypt.
	hash, err := bcrypt.GenerateFromPassword(
		[]byte("not-a-user-password"),
		12,
	)
	if err != nil {
		panic(err)
	}

	return &App{
		DB:            db,
		Origin:        origin,
		SecureCookies: secure,
		events:        events.NewRepository(db),
		dummyHash:     hash,
	}
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()

	// Stan aplikacji i połączenia z bazą.
	mux.HandleFunc("GET /api/health", a.health)

	// Logowanie i sesje.
	mux.HandleFunc("GET /api/auth/me", a.me)
	mux.HandleFunc(
		"POST /api/auth/login",
		a.mutation(a.login),
	)
	mux.HandleFunc(
		"POST /api/auth/logout",
		a.mutation(a.logout),
	)

	// Publiczne wydarzenia.
	mux.HandleFunc(
		"GET /api/events",
		a.listPublishedEvents,
	)
	mux.HandleFunc(
		"GET /api/events/{id}",
		a.getPublishedEvent,
	)

	// Wszystkie operacje administracyjne wymagają roli admin.
	mux.HandleFunc(
		"GET /api/admin/events",
		a.RequireRole([]string{"admin"}, a.listAdminEvents),
	)
	mux.HandleFunc(
		"GET /api/admin/events/{id}",
		a.RequireRole([]string{"admin"}, a.getAdminEvent),
	)

	// Zapisy wymagają dodatkowo zgodnego Origin i JSON.
	mux.HandleFunc(
		"POST /api/admin/events",
		a.mutation(
			a.RequireRole([]string{"admin"}, a.createAdminEvent),
		),
	)
	mux.HandleFunc(
		"PUT /api/admin/events/{id}",
		a.mutation(
			a.RequireRole([]string{"admin"}, a.updateAdminEvent),
		),
	)
	mux.HandleFunc(
		"DELETE /api/admin/events/{id}",
		a.mutation(
			a.RequireRole([]string{"admin"}, a.deleteAdminEvent),
		),
	)

	return mux
}

// mutation chroni operacje zmieniające dane.
func (a *App) mutation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.Origin == "" || r.Header.Get("Origin") != a.Origin {
			fail(
				w,
				http.StatusForbidden,
				"Niedozwolone źródło żądania.",
			)
			return
		}

		mediaType, _, err := mime.ParseMediaType(
			r.Header.Get("Content-Type"),
		)
		if err != nil || mediaType != "application/json" {
			fail(
				w,
				http.StatusUnsupportedMediaType,
				"Wymagany Content-Type: application/json.",
			)
			return
		}

		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}

func (a *App) dbError(w http.ResponseWriter, err error) {
	log.Printf("database error: %v", err)

	fail(
		w,
		http.StatusServiceUnavailable,
		"Baza danych jest chwilowo niedostępna.",
	)
}

// parseEventID akceptuje wyłącznie dodatnią liczbę całkowitą
// zapisaną cyframi dziesiętnymi i mieszczącą się w int64.
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

// decodeStrictJSON wymaga pojedynczego obiektu JSON.
// Odrzuca null, tablice, nieznane pola oraz dodatkowe dane
// po pierwszym obiekcie. Limit rozmiaru ustawiają handlery.
func decodeStrictJSON(r io.Reader, dst any) error {
	dec := json.NewDecoder(r)

	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return err
	}

	if len(raw) == 0 || raw[0] != '{' {
		return errors.New("expected JSON object")
	}

	if err := dec.Decode(new(any)); err != io.EOF {
		return errors.New("expected one JSON value")
	}

	// Powtórzone pola (także różniące się wielkością liter) są niejednoznaczne.
	object := json.NewDecoder(bytes.NewReader(raw))
	if _, err := object.Token(); err != nil {
		return err
	}
	seen := make(map[string]bool)
	for object.More() {
		key, err := object.Token()
		if err != nil {
			return err
		}
		name := strings.ToLower(key.(string))
		if seen[name] {
			return errors.New("duplicate JSON field")
		}
		seen[name] = true
		var value json.RawMessage
		if err := object.Decode(&value); err != nil {
			return err
		}
	}

	fields := json.NewDecoder(bytes.NewReader(raw))
	fields.DisallowUnknownFields()

	return fields.Decode(dst)
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := a.DB.PingContext(ctx); err != nil {
		a.dbError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{
		"ok": true,
	})
}
