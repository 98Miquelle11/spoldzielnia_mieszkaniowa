package events

import (
	"context"
	"strings"
	"time"
)

const (
	DefaultPageSize       = 20
	MaxPageSize           = 50
	MaxPageNumber         = 1000000
	ListDescriptionLength = 240
)

type PageOptions struct{ Page, Limit int }
type Pagination struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
}
type Page struct {
	Items      []Event    `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func (p PageOptions) Validate() error {
	if p.Page < 1 || p.Page > MaxPageNumber || p.Limit < 1 || p.Limit > MaxPageSize {
		return invalid("pagination", "nieprawidłowy numer strony lub limit")
	}
	return nil
}

func (r *Repository) page(ctx context.Context, options PageOptions, suffix string, args ...any) (Page, error) {
	if err := options.Validate(); err != nil {
		return Page{}, err
	}
	// Lista zawiera zajawkę; pełny opis jest dostępny w szczegółach i do edycji.
	listColumns := strings.Replace(columns, "description", "LEFT(description, 240) AS description", 1)
	args = append(args, options.Limit+1, (options.Page-1)*options.Limit)
	items, err := r.list(ctx, "SELECT "+listColumns+" FROM cultural_events "+suffix+" LIMIT ? OFFSET ?", args...)
	if err != nil {
		return Page{}, err
	}
	more := len(items) > options.Limit
	if more {
		items = items[:options.Limit]
	}
	return Page{Items: items, Pagination: Pagination{Page: options.Page, Limit: options.Limit, HasMore: more}}, nil
}

func (r *Repository) ListPublishedPage(ctx context.Context, period Period, options PageOptions) (Page, error) {
	now := time.Now().UTC().Truncate(time.Second)
	switch period {
	case "", Upcoming:
		return r.page(ctx, options, "WHERE is_published = ? AND ends_at >= ? ORDER BY starts_at ASC, id ASC", true, now)
	case Archived:
		return r.page(ctx, options, "WHERE is_published = ? AND ends_at < ? ORDER BY ends_at DESC, id DESC", true, now)
	default:
		return Page{}, invalid("period", "dozwolone wartości: upcoming albo past")
	}
}
func (r *Repository) ListAllPage(ctx context.Context, options PageOptions) (Page, error) {
	return r.page(ctx, options, "ORDER BY starts_at DESC, id DESC")
}
