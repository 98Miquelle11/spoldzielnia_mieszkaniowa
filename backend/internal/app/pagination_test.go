package app

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidPeriodHasPublicPolishMessage(t *testing.T) {
	a, _ := setup(t)
	w := request(a, "GET", "/api/events?period=wrong", "", nil)
	var payload errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if w.Code != 400 || payload.Error != "Nieprawidłowy filtr. Dozwolone wartości period: upcoming lub past." || strings.Contains(w.Body.String(), "invalid event input") {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
func TestPageQueryValidation(t *testing.T) {
	for _, query := range []string{"page=0", "page=-1", "page=1.5", "page=", "page=1000001", "page=999999999999999999999", "limit=0", "limit=51", "limit=abc", "limit=", "page=1&page=2", "limit=20&limit=30"} {
		a, _ := setup(t)
		w := request(a, "GET", "/api/events?"+query, "", nil)
		if w.Code != 400 {
			t.Fatalf("%s: %d %s", query, w.Code, w.Body)
		}
	}
	for _, query := range []string{"", "page=1", "page=2&limit=50"} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/events?"+query, nil)
		options, ok := eventPageOptions(w, r)
		if !ok || options.Validate() != nil {
			t.Fatalf("%s: %+v", query, options)
		}
	}
}
