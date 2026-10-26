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

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func setup(t *testing.T) (*App, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("test password"), bcrypt.MinCost)
	a := &App{DB: db, Origin: "http://localhost:3000", SecureCookies: true, dummyHash: hash}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		db.Close()
	})
	return a, mock
}

func request(a *App, method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Origin", a.Origin)
	r.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	a.Handler().ServeHTTP(w, r)
	return w
}

func TestLoginSessionLogout(t *testing.T) {
	a, mock := setup(t)
	// Spaces at the ends of a password must be preserved.
	password := "  my long password  "
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE email = ? LIMIT 1")).WithArgs("a@example.test").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).AddRow(1, "a@example.test", string(hash), "admin"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)")).WithArgs(sqlmock.AnyArg(), uint64(1), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	old := &http.Cookie{Name: cookieName, Value: strings.Repeat("b", 64)}
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sessions WHERE token_hash = ?")).WithArgs(tokenHash(old.Value)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	body, _ := json.Marshal(map[string]string{"email": " a@example.test ", "password": password})
	w := request(a, "POST", "/api/auth/login", string(body), old)
	if w.Code != http.StatusOK {
		t.Fatalf("login: %d %s", w.Code, w.Body)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies: %v", cookies)
	}
	cookie := cookies[0]
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != 28800 || len(cookie.Value) != 64 || cookie.Value == old.Value {
		t.Fatalf("unsafe cookie: %+v", cookie)
	}
	mock.ExpectQuery("SELECT u.id, u.email, u.role FROM sessions s.*JOIN users u.*s.expires_at > UTC_TIMESTAMP").WithArgs(tokenHash(cookie.Value)).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(1, "a@example.test", "admin"))
	w = request(a, "GET", "/api/auth/me", "", cookie)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"role":"admin"`) {
		t.Fatalf("me: %d %s", w.Code, w.Body)
	}
	if strings.Contains(w.Body.String(), "password") || strings.Contains(w.Body.String(), cookie.Value) {
		t.Fatal("private data exposed")
	}
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sessions WHERE token_hash = ?")).WithArgs(tokenHash(cookie.Value)).WillReturnResult(sqlmock.NewResult(0, 1))
	w = request(a, "POST", "/api/auth/logout", "{}", cookie)
	if w.Code != 200 || w.Result().Cookies()[0].MaxAge != -1 {
		t.Fatalf("logout: %d %s", w.Code, w.Body)
	}
	mock.ExpectQuery("SELECT u.id, u.email, u.role FROM sessions s.*JOIN users u.*s.expires_at > UTC_TIMESTAMP").WithArgs(tokenHash(cookie.Value)).WillReturnError(sql.ErrNoRows)
	w = request(a, "GET", "/api/auth/me", "", cookie)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"user":null`) {
		t.Fatalf("revoked session accepted: %s", w.Body)
	}
}

func TestRejectedCredentials(t *testing.T) {
	for _, known := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown", true: "wrong-password"}[known], func(t *testing.T) {
			a, mock := setup(t)
			q := mock.ExpectQuery("SELECT id, email, password_hash, role FROM users").WithArgs("a@example.test")
			if known {
				q.WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).AddRow(1, "a@example.test", string(a.dummyHash), "mieszkaniec"))
			} else {
				q.WillReturnError(sql.ErrNoRows)
			}
			w := request(a, "POST", "/api/auth/login", `{"email":"a@example.test","password":"wrong-password"}`, nil)
			if w.Code != 401 || len(w.Result().Cookies()) != 0 {
				t.Fatalf("credentials accepted: %d", w.Code)
			}
		})
	}
}

func TestCSRFAndValidation(t *testing.T) {
	a, _ := setup(t)
	for _, path := range []string{"/api/auth/login", "/api/auth/logout"} {
		for _, origin := range []string{"", "https://attacker.example"} {
			r := httptest.NewRequest("POST", path, strings.NewReader("{}"))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Origin", origin)
			w := httptest.NewRecorder()
			a.Handler().ServeHTTP(w, r)
			if w.Code != 403 {
				t.Errorf("origin %q accepted at %s", origin, path)
			}
		}
	}
	for _, body := range []string{`{`, `{}`, `{"email":"a@b.c","password":"a","extra":true}`, `{"email":"a@b.c","password":"a"} {}`, strings.Repeat("a", 5000)} {
		w := request(a, "POST", "/api/auth/login", body, nil)
		if w.Code != 400 {
			t.Errorf("invalid body accepted: %d", w.Code)
		}
	}
	r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader("email=x"))
	r.Header.Set("Origin", a.Origin)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	a.Handler().ServeHTTP(w, r)
	if w.Code != 415 {
		t.Errorf("form body accepted: %d", w.Code)
	}
	if got := request(a, "GET", "/api/auth/logout", "", nil).Code; got != 405 {
		t.Errorf("GET logout: %d", got)
	}
}

func TestExpiredSessionAndDatabaseFailure(t *testing.T) {
	a, mock := setup(t)
	cookie := &http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)}
	mock.ExpectQuery("SELECT u.id.*s.expires_at > UTC_TIMESTAMP").WithArgs(tokenHash(cookie.Value)).WillReturnError(sql.ErrNoRows)
	w := request(a, "GET", "/api/auth/me", "", cookie)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"user":null`) || w.Result().Cookies()[0].MaxAge != -1 {
		t.Fatal("expired session not cleared")
	}
	mock.ExpectQuery("SELECT u.id").WillReturnError(errors.New("secret database connection details"))
	w = request(a, "GET", "/api/auth/me", "", cookie)
	if w.Code != 503 || strings.Contains(w.Body.String(), "secret") {
		t.Fatal("database error leaked or hidden")
	}
}

func TestRoleAuthorization(t *testing.T) {
	for _, role := range []string{"", "mieszkaniec", "admin"} {
		t.Run(role, func(t *testing.T) {
			a, mock := setup(t)
			r := httptest.NewRequest("GET", "/private", nil)
			if role != "" {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)})
				mock.ExpectQuery("SELECT u.id").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow(1, "a@b.test", role))
			}
			w := httptest.NewRecorder()
			a.RequireRole([]string{"admin"}, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })(w, r)
			want := map[string]int{"": 401, "mieszkaniec": 403, "admin": 204}[role]
			if w.Code != want {
				t.Fatalf("role %q: got %d want %d", role, w.Code, want)
			}
		})
	}
}
