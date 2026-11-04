package app

import (
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPublicationStatusMustBeExplicitBoolean(t *testing.T) {
	for _, status := range []string{"", `,"is_published":null`, `,"is_published":"false"`} {
		a, mock := setup(t)
		cookie := &http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)}
		mock.ExpectQuery("SELECT u.id.*FROM sessions s.*JOIN users u").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(1, "admin@example.test", "admin"))
		body := `{"title":"Spotkanie","description":"Opis","location":"Sala","starts_at":"2026-12-01T10:00:00Z"` + status + `}`
		w := request(a, "POST", "/api/admin/events", body, cookie)
		if w.Code != 400 {
			t.Fatalf("invalid publication %s accepted: %d %s", status, w.Code, w.Body)
		}
	}
}

func TestStrictJSONRejectsDuplicateFields(t *testing.T) {
	for _, body := range []string{`{"email":"a","email":"b"}`, `{"email":"a","Email":"b"}`} {
		var value struct {
			Email string `json:"email"`
		}
		if err := decodeStrictJSON(strings.NewReader(body), &value); err == nil {
			t.Fatalf("ambiguous JSON accepted: %s", body)
		}
	}
}

func TestEveryAdminRouteRejectsGuestAndResident(t *testing.T) {
	for _, route := range []struct{ method, path, body string }{
		{"GET", "/api/admin/events", ""}, {"GET", "/api/admin/events/1", ""},
		{"POST", "/api/admin/events", "{}"}, {"PUT", "/api/admin/events/1", "{}"},
		{"DELETE", "/api/admin/events/1", "{}"},
	} {
		a, mock := setup(t)
		if w := request(a, route.method, route.path, route.body, nil); w.Code != 401 {
			t.Fatalf("guest %s %s: %d", route.method, route.path, w.Code)
		}
		cookie := &http.Cookie{Name: cookieName, Value: strings.Repeat("b", 64)}
		mock.ExpectQuery("SELECT u.id.*FROM sessions s.*JOIN users u").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(2, "resident@example.test", "mieszkaniec"))
		if w := request(a, route.method, route.path, route.body, cookie); w.Code != 403 {
			t.Fatalf("resident %s %s: %d", route.method, route.path, w.Code)
		}
	}
}
