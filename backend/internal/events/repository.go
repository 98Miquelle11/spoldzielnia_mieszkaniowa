package events

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct{ db *sql.DB }

// NewRepository reuses the connection pool opened by database.Open.
func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

// These fragments are constants. User input is always passed separately as args.
const columns = "id, title, description, starts_at, ends_at, location, is_published, created_at, updated_at"

type scanner interface{ Scan(...any) error }

func scanEvent(row scanner) (Event, error) {
	var e Event
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.StartsAt, &e.EndsAt, &e.Location, &e.IsPublished, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *Repository) list(ctx context.Context, query string, args ...any) ([]Event, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pobieranie listy wydarzeń: %w", err)
	}
	defer rows.Close()
	result := make([]Event, 0)
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("odczyt wydarzenia: %w", err)
		}
		result = append(result, e)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("odczyt listy wydarzeń: %w", err)
	}
	return result, nil
}

// ListPublished includes ongoing events in Upcoming. A past event has ended.
// No period/filter value is interpolated into SQL.
func (r *Repository) ListPublished(ctx context.Context, period Period) ([]Event, error) {
	now := time.Now().UTC().Truncate(time.Second)
	switch period {
	case "", Upcoming:
		return r.list(ctx, "SELECT "+columns+" FROM cultural_events WHERE is_published = ? AND ends_at >= ? ORDER BY starts_at ASC, id ASC", true, now)
	case Past:
		return r.list(ctx, "SELECT "+columns+" FROM cultural_events WHERE is_published = ? AND ends_at < ? ORDER BY ends_at DESC, id DESC", true, now)
	default:
		return nil, invalid("period", "dozwolone wartości: upcoming albo past")
	}
}

// ListAll is administrative: it also returns hidden events.
func (r *Repository) ListAll(ctx context.Context) ([]Event, error) {
	return r.list(ctx, "SELECT "+columns+" FROM cultural_events ORDER BY starts_at DESC, id DESC")
}

func (r *Repository) get(ctx context.Context, query string, args ...any) (Event, error) {
	e, err := scanEvent(r.db.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return Event{}, ErrNotFound
	}
	if err != nil {
		return Event{}, fmt.Errorf("pobieranie wydarzenia: %w", err)
	}
	return e, nil
}

// GetPublishedByID deliberately treats hidden and missing events identically.
func (r *Repository) GetPublishedByID(ctx context.Context, id int64) (Event, error) {
	if err := validateID(id); err != nil {
		return Event{}, err
	}
	return r.get(ctx, "SELECT "+columns+" FROM cultural_events WHERE id = ? AND is_published = ?", id, true)
}

// GetByID is administrative: it may return a hidden event.
func (r *Repository) GetByID(ctx context.Context, id int64) (Event, error) {
	if err := validateID(id); err != nil {
		return Event{}, err
	}
	return r.get(ctx, "SELECT "+columns+" FROM cultural_events WHERE id = ?", id)
}

// Create returns the generated ID. Zero-value IsPublished creates a hidden event.
func (r *Repository) Create(ctx context.Context, input Input) (int64, error) {
	input, err := normalize(input)
	if err != nil {
		return 0, err
	}
	result, err := r.db.ExecContext(ctx, `INSERT INTO cultural_events
		(title, description, starts_at, ends_at, location, is_published)
		VALUES (?, ?, ?, ?, ?, ?)`, input.Title, input.Description, input.StartsAt, input.EndsAt, input.Location, input.IsPublished)
	if err != nil {
		return 0, fmt.Errorf("dodawanie wydarzenia: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("odczyt identyfikatora wydarzenia: %w", err)
	}
	return id, nil
}

// Update locks the existing row so a no-op update is not confused with a missing
// event. MySQL normally reports changed rows, so RowsAffected()==0 is not enough.
func (r *Repository) Update(ctx context.Context, id int64, input Input) error {
	if err := validateID(id); err != nil {
		return err
	}
	input, err := normalize(input)
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("rozpoczęcie edycji wydarzenia: %w", err)
	}
	defer tx.Rollback()
	var existingID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM cultural_events WHERE id = ? FOR UPDATE", id).Scan(&existingID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("blokowanie wydarzenia do edycji: %w", err)
	}
	_, err = tx.ExecContext(ctx, `UPDATE cultural_events
		SET title = ?, description = ?, starts_at = ?, ends_at = ?, location = ?, is_published = ?
		WHERE id = ?`, input.Title, input.Description, input.StartsAt, input.EndsAt, input.Location, input.IsPublished, id)
	if err != nil {
		return fmt.Errorf("aktualizacja wydarzenia: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("zatwierdzanie edycji wydarzenia: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if err := validateID(id); err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, "DELETE FROM cultural_events WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("usuwanie wydarzenia: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sprawdzenie usunięcia wydarzenia: %w", err)
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
