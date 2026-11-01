package events

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupRepo(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		db.Close()
	})
	return NewRepository(db), mock
}

func TestListPublished_UpcomingSortAndEmptyList(t *testing.T) {
	repo, mock := setupRepo(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT " + columns + " FROM cultural_events WHERE is_published = ? AND ends_at >= ? ORDER BY starts_at ASC, id ASC")).
		WithArgs(true, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "starts_at", "ends_at", "location", "is_published", "created_at", "updated_at"}))

	items, err := repo.ListPublished(context.Background(), Upcoming)
	if err != nil {
		t.Fatalf("ListPublished upcoming: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}
}

func TestListPublished_ArchivedSort(t *testing.T) {
	repo, mock := setupRepo(t)
	now := time.Now().UTC().Truncate(time.Second)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT " + columns + " FROM cultural_events WHERE is_published = ? AND ends_at < ? ORDER BY ends_at DESC, id DESC")).
		WithArgs(true, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "starts_at", "ends_at", "location", "is_published", "created_at", "updated_at"}).
			AddRow(2, "Nowsze", "Opis", now.Add(-2*time.Hour), now.Add(-time.Hour), "Sala", true, now, now).
			AddRow(1, "Starsze", "Opis", now.Add(-4*time.Hour), now.Add(-3*time.Hour), "Sala", true, now, now))

	items, err := repo.ListPublished(context.Background(), Archived)
	if err != nil {
		t.Fatalf("ListPublished archived: %v", err)
	}
	if len(items) != 2 || items[0].ID != 2 || items[1].ID != 1 {
		t.Fatalf("unexpected order: %+v", items)
	}
}

func TestListPublished_InvalidFilter(t *testing.T) {
	repo, _ := setupRepo(t)
	_, err := repo.ListPublished(context.Background(), Period("wrong"))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetPublishedByID_HiddenOrMissingReturnsNotFound(t *testing.T) {
	repo, mock := setupRepo(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT " + columns + " FROM cultural_events WHERE id = ? AND is_published = ?")).
		WithArgs(int64(42), true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "starts_at", "ends_at", "location", "is_published", "created_at", "updated_at"}))

	_, err := repo.GetPublishedByID(context.Background(), 42)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetPublishedByID_InvalidID(t *testing.T) {
	repo, _ := setupRepo(t)
	_, err := repo.GetPublishedByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListPublished_DatabaseFailure(t *testing.T) {
	repo, mock := setupRepo(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT " + columns + " FROM cultural_events WHERE is_published = ? AND ends_at >= ? ORDER BY starts_at ASC, id ASC")).
		WithArgs(true, sqlmock.AnyArg()).
		WillReturnError(errors.New("db down"))

	_, err := repo.ListPublished(context.Background(), Upcoming)
	if err == nil {
		t.Fatal("expected database error")
	}
}
