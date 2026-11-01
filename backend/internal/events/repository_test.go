package events

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func sampleInput() Input {
	return Input{Title: "  Spotkanie  ", Description: "Opis", Location: "Sala", StartsAt: time.Date(2026, 12, 1, 10, 0, 0, 0, time.UTC), EndsAt: time.Date(2026, 12, 1, 11, 0, 0, 0, time.UTC)}
}
func TestPublishedBetweenOverlapQuery(t *testing.T) {
	repo, mock := setupRepo(t)
	from := time.Date(2026, 12, 1, 0, 0, 0, 0, time.FixedZone("offset", 3600))
	to := from.Add(24 * time.Hour)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT "+columns+" FROM cultural_events WHERE is_published = ? AND starts_at < ? AND ends_at >= ? ORDER BY starts_at ASC, id ASC")).WithArgs(true, to.UTC(), from.UTC()).WillReturnRows(sqlmock.NewRows(strings.Split(columns, ", ")))
	result, err := repo.ListPublishedBetween(context.Background(), from, to)
	if err != nil || result == nil || len(result) != 0 {
		t.Fatalf("%v %v", result, err)
	}
	for _, pair := range [][2]time.Time{{from, from}, {to, from}, {time.Time{}, to}} {
		if _, err = repo.ListPublishedBetween(context.Background(), pair[0], pair[1]); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("invalid interval: %v", err)
		}
	}
}
func TestCreateParameterizedAndValidated(t *testing.T) {
	repo, mock := setupRepo(t)
	input := sampleInput()
	input.Title = "'; DROP TABLE users; --"
	mock.ExpectExec("INSERT INTO cultural_events").WithArgs(input.Title, input.Description, input.StartsAt, input.EndsAt, input.Location, false).WillReturnResult(sqlmock.NewResult(8, 1))
	id, err := repo.Create(context.Background(), input)
	if err != nil || id != 8 {
		t.Fatalf("%d %v", id, err)
	}
	input.Title = ""
	if _, err = repo.Create(context.Background(), input); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("missing title accepted: %v", err)
	}
}
func TestUpdateUnchangedRowAndMissing(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "unchanged", true: "missing"}[missing], func(t *testing.T) {
			repo, mock := setupRepo(t)
			mock.ExpectBegin()
			query := mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM cultural_events WHERE id = ? FOR UPDATE")).WithArgs(int64(9))
			if missing {
				query.WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			} else {
				query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(9))
				input := sampleInput()
				mock.ExpectExec("UPDATE cultural_events").WithArgs("Spotkanie", input.Description, input.StartsAt, input.EndsAt, input.Location, false, int64(9)).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			}
			err := repo.Update(context.Background(), 9, sampleInput())
			if missing && !errors.Is(err, ErrNotFound) || !missing && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
func TestDeleteResultAndFailure(t *testing.T) {
	for _, count := range []int64{0, 1} {
		repo, mock := setupRepo(t)
		mock.ExpectExec("DELETE FROM cultural_events WHERE id =").WithArgs(int64(9)).WillReturnResult(sqlmock.NewResult(0, count))
		err := repo.Delete(context.Background(), 9)
		if count == 0 && !errors.Is(err, ErrNotFound) || count == 1 && err != nil {
			t.Fatalf("%v", err)
		}
	}
	repo, mock := setupRepo(t)
	mock.ExpectExec("DELETE FROM cultural_events").WillReturnError(errors.New("offline"))
	if err := repo.Delete(context.Background(), 9); err == nil {
		t.Fatal("SQL error hidden")
	}
}
func TestInputDateAndTextValidation(t *testing.T) {
	input := sampleInput()
	input.EndsAt = input.StartsAt.Add(-time.Second)
	if ValidateInput(input)["ends_at"] == "" {
		t.Fatal("reversed dates accepted")
	}
	input.EndsAt = input.StartsAt
	if len(ValidateInput(input)) != 0 {
		t.Fatal("equal dates rejected")
	}
	input.Title = strings.Repeat("ą", 201)
	if ValidateInput(input)["title"] == "" {
		t.Fatal("title limit not enforced")
	}
	input = sampleInput()
	input.StartsAt = time.Date(999, 1, 1, 0, 0, 0, 0, time.UTC)
	if ValidateInput(input)["starts_at"] == "" {
		t.Fatal("MySQL range not enforced")
	}
}
