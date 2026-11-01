package events

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxTitleLength       = 200
	MaxLocationLength    = 200
	MaxDescriptionLength = 5000
)

var (
	ErrNotFound = errors.New("event not found")
)

type Event struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	IsPublished bool       `json:"is_published"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpsertEventInput reprezentuje dane przesyłane z formularza admina (POST / PUT).
type UpsertEventInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      *string `json:"ends_at,omitempty"`
	IsPublished bool    `json:"is_published"`
}

// ValidatedEvent przechowuje zwalidowane i sparsowane dane gotowe do zapisu w repozytorium.
type ValidatedEvent struct {
	Title       string
	Description string
	Location    string
	StartsAt    time.Time
	EndsAt      *time.Time
	IsPublished bool
}

// Validate sprawdza wymagane pola, długości tekstów oraz poprawność terminów RFC 3339.
// Zwraca mapę błędów dla poszczególnych pól formularza.
func (in UpsertEventInput) Validate() (ValidatedEvent, map[string]string) {
	fieldErrors := make(map[string]string)

	title := strings.TrimSpace(in.Title)
	if title == "" {
		fieldErrors["title"] = "Tytuł wydarzenia jest wymagany."
	} else if utf8.RuneCountInString(title) > MaxTitleLength {
		fieldErrors["title"] = "Tytuł nie może przekraczać 200 znaków."
	}

	description := strings.TrimSpace(in.Description)
	if description == "" {
		fieldErrors["description"] = "Opis wydarzenia jest wymagany."
	} else if utf8.RuneCountInString(description) > MaxDescriptionLength {
		fieldErrors["description"] = "Opis nie może przekraczać 5000 znaków."
	}

	location := strings.TrimSpace(in.Location)
	if location == "" {
		fieldErrors["location"] = "Miejsce wydarzenia jest wymagane."
	} else if utf8.RuneCountInString(location) > MaxLocationLength {
		fieldErrors["location"] = "Miejsce nie może przekraczać 200 znaków."
	}

	var startsAt time.Time
	startsAtStr := strings.TrimSpace(in.StartsAt)
	if startsAtStr == "" {
		fieldErrors["starts_at"] = "Termin rozpoczęcia jest wymagany."
	} else {
		parsedStart, err := time.Parse(time.RFC3339, startsAtStr)
		if err != nil {
			fieldErrors["starts_at"] = "Nieprawidłowy format daty rozpoczęcia (wymagany RFC 3339)."
		} else {
			startsAt = parsedStart.UTC()
		}
	}

	var endsAt *time.Time
	if in.EndsAt != nil && strings.TrimSpace(*in.EndsAt) != "" {
		parsedEnd, err := time.Parse(time.RFC3339, strings.TrimSpace(*in.EndsAt))
		if err != nil {
			fieldErrors["ends_at"] = "Nieprawidłowy format daty zakończenia (wymagany RFC 3339)."
		} else {
			utcEnd := parsedEnd.UTC()
			if !startsAt.IsZero() && utcEnd.Before(startsAt) {
				fieldErrors["ends_at"] = "Termin zakończenia nie może być wcześniejszy niż termin rozpoczęcia."
			} else {
				endsAt = &utcEnd
			}
		}
	}

	if len(fieldErrors) > 0 {
		return ValidatedEvent{}, fieldErrors
	}

	return ValidatedEvent{
		Title:       title,
		Description: description,
		Location:    location,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		IsPublished: in.IsPublished,
	}, nil
}