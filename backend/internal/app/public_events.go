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
	if period != "" && period != events.Upcoming && period != events.Archived {
		fail(w, http.StatusBadRequest, "Nieprawidłowy filtr. Dozwolone wartości period: upcoming lub past.")
		return
	}
	options, ok := eventPageOptions(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := a.events.ListPublishedPage(ctx, period, options)
	if err != nil {
		if errors.Is(err, events.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Nieprawidłowe parametry listy wydarzeń."})
			return
		}
		a.dbError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
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
