package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"mirmilowo/backend/internal/events"
)

func (a *App) listPublishedEvents(w http.ResponseWriter, r *http.Request) {
	period := events.Period(r.URL.Query().Get("period"))
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := a.events.ListPublished(ctx, period)
	if err != nil {
		if errors.Is(err, events.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		a.dbError(w, err)
		return
	}

	if items == nil {
		items = []events.Event{}
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (a *App) getPublishedEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseEventID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Nieprawidłowy identyfikator wydarzenia."})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	event, err := a.events.GetPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, events.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "Nie znaleziono wydarzenia."})
			return
		}
		a.dbError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, event)
}
