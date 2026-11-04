package events

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPageQueriesBoundedAndSummarized(t *testing.T) {
	for _, period := range []Period{Upcoming, Archived} {
		repo, mock := setupRepo(t)
		condition, order := "ends_at >= ?", "starts_at ASC, id ASC"
		if period == Archived {
			condition, order = "ends_at < ?", "ends_at DESC, id DESC"
		}
		query := "SELECT " + strings.Replace(columns, "description", "LEFT(description, 240) AS description", 1) + " FROM cultural_events WHERE is_published = ? AND " + condition + " ORDER BY " + order + " LIMIT ? OFFSET ?"
		rows := sqlmock.NewRows(strings.Split(columns, ", "))
		now := time.Now().UTC()
		for i := 1; i <= 3; i++ {
			rows.AddRow(i, "Tytuł", "Opis", now, now, "Sala", true, now, now)
		}
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(true, sqlmock.AnyArg(), 3, 2).WillReturnRows(rows)
		got, err := repo.ListPublishedPage(context.Background(), period, PageOptions{Page: 2, Limit: 2})
		if err != nil || len(got.Items) != 2 || !got.Pagination.HasMore || got.Pagination.Page != 2 {
			t.Fatalf("%+v %v", got, err)
		}
	}
}
func TestAdminPageContainsHiddenAndEmptyPage(t *testing.T) {
	repo, mock := setupRepo(t)
	query := "SELECT " + strings.Replace(columns, "description", "LEFT(description, 240) AS description", 1) + " FROM cultural_events ORDER BY starts_at DESC, id DESC LIMIT ? OFFSET ?"
	rows := sqlmock.NewRows(strings.Split(columns, ", "))
	now := time.Now().UTC()
	rows.AddRow(1, "Ukryte", "Opis", now, now, "Sala", false, now, now)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(21, 0).WillReturnRows(rows)
	got, err := repo.ListAllPage(context.Background(), PageOptions{Page: 1, Limit: 20})
	if err != nil || len(got.Items) != 1 || got.Items[0].IsPublished || got.Pagination.HasMore {
		t.Fatalf("%+v %v", got, err)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(21, 20).WillReturnRows(sqlmock.NewRows(strings.Split(columns, ", ")))
	got, err = repo.ListAllPage(context.Background(), PageOptions{Page: 2, Limit: 20})
	if err != nil || got.Items == nil || len(got.Items) != 0 || got.Pagination.HasMore {
		t.Fatalf("%+v %v", got, err)
	}
}
func TestInvalidPageNeverQueriesDatabase(t *testing.T) {
	repo, _ := setupRepo(t)
	for _, options := range []PageOptions{{Page: 0, Limit: 20}, {Page: 1, Limit: 0}, {Page: 1, Limit: 51}, {Page: MaxPageNumber + 1, Limit: 20}} {
		if _, err := repo.ListAllPage(context.Background(), options); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("accepted %+v: %v", options, err)
		}
	}
}
