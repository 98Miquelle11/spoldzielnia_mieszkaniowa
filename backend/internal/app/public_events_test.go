package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPublicListContract(t *testing.T) {
	for _, period := range []string{"", "upcoming", "past"} {
		t.Run("period="+period, func(t *testing.T) {
			a, mock := setup(t)
			condition, order := "ends_at >= ?", "starts_at ASC, id ASC"
			if period == "past" {
				condition, order = "ends_at < ?", "ends_at DESC, id DESC"
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT "+eventsColumnsForTest+" FROM cultural_events WHERE is_published = ? AND "+condition+" ORDER BY "+order)).WithArgs(true, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows(strings.Split(eventsColumnsForTest, ", ")))
			w := request(a, "GET", "/api/events?period="+period, "", nil)
			if w.Code != 200 || strings.TrimSpace(w.Body.String()) != `{"items":[]}` {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
		})
	}
}
func TestPublicListInvalidPeriodAndDatabaseFailure(t *testing.T) {
	a, mock := setup(t)
	if w := request(a, "GET", "/api/events?period=invalid", "", nil); w.Code != 400 {
		t.Fatalf("invalid period: %d", w.Code)
	}
	mock.ExpectQuery("SELECT .*FROM cultural_events").WillReturnError(errors.New("private SQL details"))
	w := request(a, "GET", "/api/events", "", nil)
	if w.Code != 503 || strings.Contains(w.Body.String(), "private SQL") {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
func TestPublicDetailContractAndUTC(t *testing.T) {
	a, mock := setup(t)
	local := time.Date(2026, 12, 1, 10, 0, 0, 0, time.FixedZone("example", 3600))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT "+eventsColumnsForTest+" FROM cultural_events WHERE id = ? AND is_published = ?")).WithArgs(int64(7), true).WillReturnRows(sqlmock.NewRows(strings.Split(eventsColumnsForTest, ", ")).AddRow(7, "Tytuł", "Opis", local, local.Add(time.Hour), "Sala", true, local, local))
	w := request(a, "GET", "/api/events/7", "", nil)
	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || result["id"] != float64(7) || result["starts_at"] != "2026-12-01T09:00:00Z" || result["event"] != nil {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
func TestPublicDetailHiddenMissingAndMalformedID(t *testing.T) {
	a, mock := setup(t)
	for _, id := range []string{"0", "-1", "abc", "9223372036854775808"} {
		if w := request(a, "GET", "/api/events/"+id, "", nil); w.Code != 400 {
			t.Fatalf("id %s: %d", id, w.Code)
		}
	}
	mock.ExpectQuery("SELECT .*FROM cultural_events WHERE id = .* AND is_published =").WithArgs(int64(42), true).WillReturnError(sql.ErrNoRows)
	if w := request(a, "GET", "/api/events/42", "", nil); w.Code != 404 {
		t.Fatalf("hidden or missing: %d", w.Code)
	}
	mock.ExpectQuery("SELECT .*FROM cultural_events").WillReturnError(errors.New("connection details"))
	if w := request(a, "GET", "/api/events/43", "", nil); w.Code != 503 || strings.Contains(w.Body.String(), "connection details") {
		t.Fatalf("database error: %d %s", w.Code, w.Body)
	}
}
func TestAllAdminRoutesRejectGuestsAndResidents(t *testing.T) {
	for _, route := range []struct{ method, path string }{
		{"GET", "/api/admin/events"}, {"GET", "/api/admin/events/1"}, {"POST", "/api/admin/events"}, {"PUT", "/api/admin/events/1"}, {"DELETE", "/api/admin/events/1"},
	} {
		for _, resident := range []bool{false, true} {
			t.Run(route.method+route.path+map[bool]string{false: " guest", true: " resident"}[resident], func(t *testing.T) {
				a, mock := setup(t)
				var cookie *http.Cookie
				want := 401
				if resident {
					want = 403
					cookie = &http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)}
					mock.ExpectQuery("SELECT u.id").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(2, "r@example.test", "mieszkaniec"))
				}
				w := request(a, route.method, route.path, "{}", cookie)
				if w.Code != want {
					t.Fatalf("%d %s", w.Code, w.Body)
				}
			})
		}
	}
}
func TestAdminValidationHasFieldsAndRejectsTrailingJSON(t *testing.T) {
	for _, body := range []string{`{}`, `{} {}`} {
		a, mock := setup(t)
		mock.ExpectQuery("SELECT u.id").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(1, "a@example.test", "admin"))
		w := request(a, "POST", "/api/admin/events", body, &http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)})
		if w.Code != 400 {
			t.Fatalf("%d %s", w.Code, w.Body)
		}
		if body == "{}" && !strings.Contains(w.Body.String(), `"fields"`) {
			t.Fatal("missing field errors")
		}
	}
}

func TestAdminDetailReturnsUnwrappedHiddenEvent(t *testing.T) {
	a, mock := setup(t)
	cookie := &http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)}
	mock.ExpectQuery("SELECT u.id").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(1, "a@example.test", "admin"))
	now := time.Date(2026, 12, 1, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT " + eventsColumnsForTest + " FROM cultural_events WHERE id = ?")).WithArgs(int64(7)).WillReturnRows(sqlmock.NewRows(strings.Split(eventsColumnsForTest, ", ")).AddRow(7, "Ukryte", "Opis", now, now.Add(time.Hour), "Sala", false, now, now))
	w := request(a, "GET", "/api/admin/events/7", "", cookie)
	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || result["id"] != float64(7) || result["is_published"] != false || result["event"] != nil {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
func TestEveryAdminMutationChecksOriginAndJSON(t *testing.T) {
	for _, route := range []struct{ method, path string }{{"POST", "/api/admin/events"}, {"PUT", "/api/admin/events/7"}, {"DELETE", "/api/admin/events/7"}} {
		a, _ := setup(t)
		for _, badOrigin := range []string{"", "https://attacker.example"} {
			r := httptest.NewRequest(route.method, route.path, strings.NewReader("{}"))
			r.Header.Set("Origin", badOrigin)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			a.Handler().ServeHTTP(w, r)
			if w.Code != 403 {
				t.Fatalf("%s Origin %q: %d", route.method, badOrigin, w.Code)
			}
		}
		r := httptest.NewRequest(route.method, route.path, strings.NewReader("{}"))
		r.Header.Set("Origin", a.Origin)
		w := httptest.NewRecorder()
		a.Handler().ServeHTTP(w, r)
		if w.Code != 415 {
			t.Fatalf("%s accepted without JSON: %d", route.method, w.Code)
		}
	}
}
