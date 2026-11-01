package app

type AdminEventsRepository interface {
    ListAll(ctx context.Context) ([]events.Event, error)
    GetByIDAdmin(ctx context.Context, id int64) (events.Event, error)
    Create(ctx context.Context, in events.ValidatedEvent) (int64, error)
    Update(ctx context.Context, id int64, in events.ValidatedEvent) error
    Delete(ctx context.Context, id int64) error
}

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"spoldzielnia_mieszkaniowa/backend/internal/events"
)

// Struktury odpowiedzi JSON
type eventsListResponse struct {
	Items []events.Event `json:"items"`
}

type createdIDResponse struct {
	ID int64 `json:"id"`
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Middleware ochrony żądań zmieniających dane (CSRF: Origin == APP_ORIGIN + Content-Type: application/json)
func (a *App) requireJSONWriteProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			origin := r.Header.Get("Origin")
			if origin == "" || origin != a.Config.AppOrigin {
				writeJSON(w, http.StatusForbidden, errorResponse{
					Error: "Niedozwolone źródło żądania (błąd Origin).",
				})
				return
			}

			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				writeJSON(w, http.StatusUnsupportedMediaType, errorResponse{
					Error: "Wymagany nagłówek Content-Type: application/json.",
				})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Rejestracja tras administracyjnych w routerze aplikacji
func (a *App) registerAdminEventRoutes(mux *http.ServeMux) {
	adminRoleMiddleware := a.RequireRole("admin")

	// Operacje odczytu (chronione rolą admin)
	mux.Handle("GET /api/admin/events", adminRoleMiddleware(http.HandlerFunc(a.handleAdminListEvents)))
	mux.Handle("GET /api/admin/events/{id}", adminRoleMiddleware(http.HandlerFunc(a.handleAdminGetEvent)))

	// Operacje zapisu (chronione rolą admin + ochroną CSRF/JSON)
	mux.Handle("POST /api/admin/events", adminRoleMiddleware(a.requireJSONWriteProtection(http.HandlerFunc(a.handleAdminCreateEvent))))
	mux.Handle("PUT /api/admin/events/{id}", adminRoleMiddleware(a.requireJSONWriteProtection(http.HandlerFunc(a.handleAdminUpdateEvent))))
	mux.Handle("DELETE /api/admin/events/{id}", adminRoleMiddleware(a.requireJSONWriteProtection(http.HandlerFunc(a.handleAdminDeleteEvent))))
}

// GET /api/admin/events
func (a *App) handleAdminListEvents(w http.ResponseWriter, r *http.Request) {
	items, err := a.EventsRepo.ListAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "Nie udało się pobrać listy wydarzeń.",
		})
		return
	}
	if items == nil {
		items = []events.Event{}
	}
	writeJSON(w, http.StatusOK, eventsListResponse{Items: items})
}

// GET /api/admin/events/{id}
func (a *App) handleAdminGetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseEventID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Nieprawidłowy identyfikator wydarzenia.",
		})
		return
	}

	ev, err := a.EventsRepo.GetByIDAdmin(r.Context(), id)
	if err != nil {
		if errors.Is(err, events.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "Wydarzenie nie istnieje.",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "Nie udało się pobrać szczegółów wydarzenia.",
		})
		return
	}

	writeJSON(w, http.StatusOK, ev)
}

// POST /api/admin/events
func (a *App) handleAdminCreateEvent(w http.ResponseWriter, r *http.Request) {
	var input events.UpsertEventInput
	if err := decodeStrictJSON(r.Body, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Nieprawidłowy format danych JSON.",
		})
		return
	}

	validated, fieldErrors := input.Validate()
	if len(fieldErrors) > 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error:  "Formularz zawiera błędy walidacji.",
			Fields: fieldErrors,
		})
		return
	}

	newID, err := a.EventsRepo.Create(r.Context(), validated)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "Nie udało się zapisać nowego wydarzenia.",
		})
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/admin/events/%d", newID))
	writeJSON(w, http.StatusCreated, createdIDResponse{ID: newID})
}

// PUT /api/admin/events/{id}
func (a *App) handleAdminUpdateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseEventID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Nieprawidłowy identyfikator wydarzenia.",
		})
		return
	}

	var input events.UpsertEventInput
	if err := decodeStrictJSON(r.Body, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Nieprawidłowy format danych JSON.",
		})
		return
	}

	validated, fieldErrors := input.Validate()
	if len(fieldErrors) > 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error:  "Formularz zawiera błędy walidacji.",
			Fields: fieldErrors,
		})
		return
	}

	if err := a.EventsRepo.Update(r.Context(), id, validated); err != nil {
		if errors.Is(err, events.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "Wydarzenie nie istnieje.",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "Nie udało się zaktualizować wydarzenia.",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/admin/events/{id}
func (a *App) handleAdminDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseEventID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Nieprawidłowy identyfikator wydarzenia.",
		})
		return
	}

	// Zgodnie z ustaleniami Dew 3. i 6., DELETE przyjmuje Content-Type: application/json i treść {}
	var emptyPayload map[string]any
	if err := decodeStrictJSON(r.Body, &emptyPayload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "Wymagane poprawne ciało JSON (np. {}).",
		})
		return
	}

	if err := a.EventsRepo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, events.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "Wydarzenie nie istnieje.",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "Nie udało się usunąć wydarzenia.",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Funkcje pomocnicze
func parseEventID(raw string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func decodeStrictJSON(r io.Reader, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected extra JSON data")
	}
	return nil
}