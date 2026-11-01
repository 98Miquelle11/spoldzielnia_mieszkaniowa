package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"mirmilowo/backend/internal/events"
)

func (a *App) listAdminEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := a.events.ListAll(ctx)
	if err != nil {
		a.dbError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]any{"events": result})
}

func (a *App) getAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := a.events.GetByID(ctx, id)
	if err != nil {
		a.eventError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]any{"event": result})
}

func (a *App) createAdminEvent(w http.ResponseWriter, r *http.Request) {
	var input events.Input
	if !decodeAdminEventInput(w, r, &input) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id, err := a.events.Create(ctx, input)
	if err != nil {
		a.eventError(w, err)
		return
	}
	result, err := a.events.GetByID(ctx, id)
	if err != nil {
		a.eventError(w, err)
		return
	}
	respond(w, http.StatusCreated, map[string]any{"event": result})
}

func (a *App) updateAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	var input events.Input
	if !decodeAdminEventInput(w, r, &input) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := a.events.Update(ctx, id, input); err != nil {
		a.eventError(w, err)
		return
	}
	result, err := a.events.GetByID(ctx, id)
	if err != nil {
		a.eventError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]any{"event": result})
}

func (a *App) setAdminEventPublication(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	var input struct {
		IsPublished *bool `json:"is_published"`
	}
	if !decodeAdminEventInput(w, r, &input) {
		return
	}
	if input.IsPublished == nil {
		fail(w, http.StatusBadRequest, "Podaj status publikacji wydarzenia.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := a.events.SetPublished(ctx, id, *input.IsPublished); err != nil {
		a.eventError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]bool{"is_published": *input.IsPublished})
}

func (a *App) deleteAdminEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := adminEventID(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := a.events.Delete(ctx, id); err != nil {
		a.eventError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]bool{"ok": true})
}

func adminEventID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		fail(w, http.StatusBadRequest, "Nieprawidłowy identyfikator wydarzenia.")
		return 0, false
	}
	return id, true
}

func decodeAdminEventInput(w http.ResponseWriter, r *http.Request, input any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 128*1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(input); err != nil {
		fail(w, http.StatusBadRequest, "Nieprawidłowe dane formularza.")
		return false
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		fail(w, http.StatusBadRequest, "Nieprawidłowe dane formularza.")
		return false
	}
	return true
}

func (a *App) eventError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, events.ErrNotFound):
		fail(w, http.StatusNotFound, "Nie znaleziono wydarzenia.")
	case errors.Is(err, events.ErrInvalidInput):
		fail(w, http.StatusBadRequest, "Nieprawidłowe dane wydarzenia: "+err.Error())
	default:
		a.dbError(w, err)
	}
}
