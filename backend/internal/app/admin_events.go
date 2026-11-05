package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"mirmilowo/backend/internal/events"
)

func (a *App) listAdminEvents(w http.ResponseWriter, r *http.Request) {
	options, ok := eventPageOptions(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := a.events.ListAllPage(ctx, options)
	if err != nil {
		a.dbError(w, r, err)
		return
	}
	writeJSON(w, 200, result)
}
func (a *App) getAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	event, err := a.events.GetByID(ctx, id)
	if err != nil {
		a.eventError(w, r, err)
		return
	}
	writeJSON(w, 200, event)
}
func (a *App) createAdminEvent(w http.ResponseWriter, r *http.Request) {
	input, ok := decodeAdminEventInput(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id, err := a.events.Create(ctx, input)
	if err != nil {
		a.eventError(w, r, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/admin/events/%d", id))
	writeJSON(w, 201, map[string]int64{"id": id})
}
func (a *App) updateAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	input, ok := decodeAdminEventInput(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := a.events.Update(ctx, id, input); err != nil {
		a.eventError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (a *App) deleteAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	var input struct{}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := decodeStrictJSON(r.Body, &input); err != nil {
		fail(w, 400, "Wyślij pusty obiekt JSON: {}.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := a.events.Delete(ctx, id); err != nil {
		a.eventError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func adminEventID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := parseEventID(r.PathValue("id"))
	if err != nil {
		fail(w, 400, "Nieprawidłowy identyfikator wydarzenia.")
		return 0, false
	}
	return id, true
}
func decodeAdminEventInput(w http.ResponseWriter, r *http.Request) (events.Input, bool) {
	// Wskaźnik odróżnia jawne false od brakującej wartości lub null.
	var payload struct {
		events.Input
		Publication *bool `json:"is_published"`
	}
	var input events.Input
	r.Body = http.MaxBytesReader(w, r.Body, 128*1024)
	if err := decodeStrictJSON(r.Body, &payload); err != nil {
		fail(w, 400, "Nieprawidłowe dane formularza. Daty muszą mieć format RFC 3339.")
		return input, false
	}
	input = payload.Input
	if payload.Publication == nil {
		writeJSON(w, 400, errorResponse{Error: "Popraw zaznaczone pola.", Fields: map[string]string{"is_published": "Podaj status publikacji jako true lub false."}})
		return input, false
	}
	input.IsPublished = *payload.Publication
	if fields := events.ValidateInput(input); len(fields) > 0 {
		writeJSON(w, 400, errorResponse{Error: "Popraw zaznaczone pola.", Fields: fields})
		return input, false
	}
	return input, true
}
func (a *App) eventError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, events.ErrNotFound):
		fail(w, 404, "Nie znaleziono wydarzenia.")
	case errors.Is(err, events.ErrInvalidInput):
		fail(w, 400, "Nieprawidłowe dane wydarzenia.")
	default:
		a.dbError(w, r, err)
	}
}
