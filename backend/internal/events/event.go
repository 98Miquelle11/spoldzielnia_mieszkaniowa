package events

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrNotFound     = errors.New("nie znaleziono wydarzenia")
	ErrInvalidInput = errors.New("nieprawidłowe dane wydarzenia")
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Location    string    `json:"location"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Input struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Location    string    `json:"location"`
	IsPublished bool      `json:"is_published"`
}

type Period string

const (
	Upcoming Period = "upcoming"
	Past     Period = "past"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string { return fmt.Sprintf("%s: %s", e.Field, e.Message) }
func (e *ValidationError) Unwrap() error { return ErrInvalidInput }

func invalid(field, message string) error { return &ValidationError{Field: field, Message: message} }

func normalize(input Input) (Input, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Location = strings.TrimSpace(input.Location)
	for _, field := range []struct {
		name, value string
		max         int
	}{
		{"title", input.Title, 200}, {"location", input.Location, 255},
	} {
		if !utf8.ValidString(field.value) || field.value == "" || utf8.RuneCountInString(field.value) > field.max {
			return input, invalid(field.name, fmt.Sprintf("wymagany poprawny tekst UTF-8, maksymalnie %d znaków", field.max))
		}
	}
	if !utf8.ValidString(input.Description) || input.Description == "" || len(input.Description) > 65535 {
		return input, invalid("description", "opis jest wymagany i może mieć maksymalnie 65535 bajtów UTF-8")
	}
	if input.StartsAt.IsZero() {
		return input, invalid("starts_at", "podaj termin rozpoczęcia")
	}
	if input.EndsAt.IsZero() {
		return input, invalid("ends_at", "podaj termin zakończenia")
	}
	input.StartsAt = input.StartsAt.UTC()
	input.EndsAt = input.EndsAt.UTC()
	for _, field := range []struct {
		name  string
		value time.Time
	}{{"starts_at", input.StartsAt}, {"ends_at", input.EndsAt}} {
		if field.value.Year() < 1000 || field.value.Year() > 9999 {
			return input, invalid(field.name, "rok musi mieścić się w zakresie 1000–9999")
		}
	}
	if input.EndsAt.Before(input.StartsAt) {
		return input, invalid("ends_at", "zakończenie nie może poprzedzać rozpoczęcia")
	}
	input.StartsAt = input.StartsAt.Truncate(time.Second)
	input.EndsAt = input.EndsAt.Truncate(time.Second)
	return input, nil
}

func validateID(id int64) error {
	if id <= 0 {
		return invalid("id", "identyfikator musi być dodatni")
	}
	return nil
}
