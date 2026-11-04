package app

import (
	"mirmilowo/backend/internal/events"
	"net/http"
	"strconv"
)

func eventPageOptions(w http.ResponseWriter, r *http.Request) (events.PageOptions, bool) {
	options := events.PageOptions{Page: 1, Limit: events.DefaultPageSize}
	query := r.URL.Query()
	for _, field := range []struct {
		name string
		dst  *int
		max  int
	}{
		{"page", &options.Page, events.MaxPageNumber}, {"limit", &options.Limit, events.MaxPageSize},
	} {
		values, present := query[field.name]
		if !present {
			continue
		}
		if len(values) != 1 {
			fail(w, 400, "Parametry page i limit podaj tylko raz.")
			return options, false
		}
		value := values[0]
		valid := len(value) > 0
		for _, ch := range value {
			if ch < '0' || ch > '9' {
				valid = false
			}
		}
		n, err := strconv.Atoi(value)
		if !valid || err != nil || n < 1 || n > field.max {
			fail(w, 400, "Nieprawidłowa paginacja: page od 1 do 1000000, limit od 1 do 50.")
			return options, false
		}
		*field.dst = n
	}
	return options, true
}
