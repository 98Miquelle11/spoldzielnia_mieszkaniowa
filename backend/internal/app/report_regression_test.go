package app

import (
	"bytes"
	"errors"
	"log"
	"mirmilowo/backend/internal/events"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStrictJSONUnicodeDuplicates(t *testing.T) {
	for _, body := range []string{
		`{"password":"bad","paſſword":"good"}`,
		`{"paſſword":"good","password":"bad"}`,
		`{"is_published":true,"iſ_publiſhed":false}`,
		`{"key":"one","Key":"two"}`,
		`{"password":"bad","pa\u017f\u017fword":"good"}`,
	} {
		t.Run(body, func(t *testing.T) {
			var dst struct {
				Password  string `json:"password"`
				Published bool   `json:"is_published"`
				Key       string `json:"key"`
			}
			if err := decodeStrictJSON(strings.NewReader(body), &dst); err == nil {
				t.Fatal("accepted duplicate fields")
			}
		})
	}
	var dst struct {
		Password string `json:"password"`
		Key      string `json:"key"`
	}
	if err := decodeStrictJSON(strings.NewReader(`{"password":"ok","key":"value"}`), &dst); err != nil {
		t.Fatal(err)
	}
}

func TestRequestLoggingCorrelatesDatabaseFailure(t *testing.T) {
	var output bytes.Buffer
	old := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(old) })
	a := &App{}
	h := requestLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.dbError(w, r, errors.New("database offline"))
	}))
	ids := map[string]bool{}
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/events?token=secret", nil)
		r.Header.Set("X-Request-ID", "untrusted")
		h.ServeHTTP(w, r)
		id := w.Header().Get("X-Request-ID")
		if id == "" || id == "untrusted" || ids[id] {
			t.Fatalf("invalid request id: %q", id)
		}
		ids[id] = true
		if w.Code != 503 || strings.Contains(w.Body.String(), "database offline") {
			t.Fatalf("response: %d %s", w.Code, w.Body)
		}
		if strings.Count(output.String(), "request_id="+id) != 2 {
			t.Fatal("database and access log must share request id")
		}
	}
	for _, field := range []string{`method=GET`, `path="/api/events"`, `status=503`, `duration=`} {
		if !strings.Contains(output.String(), field) {
			t.Fatalf("missing %s", field)
		}
	}
	if strings.Contains(output.String(), "secret") {
		t.Fatal("query string leaked")
	}
}

func TestEventErrorMapping(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
	}{
		{events.ErrNotFound, 404}, {events.ErrInvalidInput, 400}, {errors.New("private SQL"), 503},
	} {
		w := httptest.NewRecorder()
		(&App{}).eventError(w, httptest.NewRequest("GET", "/api/admin/events/1", nil), tc.err)
		if w.Code != tc.status || strings.Contains(w.Body.String(), "private SQL") {
			t.Fatalf("response %d: %s", w.Code, w.Body)
		}
	}
}

func TestRequestLoggingStatus(t *testing.T) {
	for _, status := range []int{200, 204, 400, 401, 404, 405} {
		w := httptest.NewRecorder()
		requestLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status) })).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		if w.Code != status || w.Header().Get("X-Request-ID") == "" {
			t.Fatalf("lost status/header: %d", status)
		}
	}
}
