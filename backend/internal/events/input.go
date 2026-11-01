package events

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type Period string

const (
	Upcoming Period = "upcoming"
	Archived Period = "past"
)

type Input struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	IsPublished bool      `json:"is_published"`
}

func ValidateInput(input Input) map[string]string {
	fields := make(map[string]string)
	for _, field := range []struct {
		key, value string
		limit      int
	}{
		{"title", input.Title, MaxTitleLength},
		{"description", input.Description, MaxDescriptionLength},
		{"location", input.Location, MaxLocationLength},
	} {
		value := strings.TrimSpace(field.value)
		if value == "" {
			fields[field.key] = "Pole jest wymagane."
		} else if utf8.RuneCountInString(value) > field.limit {
			fields[field.key] = fmt.Sprintf("Maksymalnie %d znaków.", field.limit)
		}
	}
	for _, field := range []struct {
		key   string
		value time.Time
	}{{"starts_at", input.StartsAt}, {"ends_at", input.EndsAt}} {
		if field.value.IsZero() {
			fields[field.key] = "Termin jest wymagany."
		} else if field.value.UTC().Year() < 1000 || field.value.UTC().Year() > 9999 {
			fields[field.key] = "Rok musi należeć do zakresu 1000–9999 (MySQL DATETIME)."
		}
	}
	if !input.StartsAt.IsZero() && !input.EndsAt.IsZero() && input.EndsAt.Before(input.StartsAt) {
		fields["ends_at"] = "Zakończenie nie może być wcześniejsze niż rozpoczęcie."
	}
	return fields
}
func validateID(id int64) error {
	if id <= 0 {
		return invalid("id", "musi być dodatnią liczbą całkowitą")
	}
	return nil
}
func normalize(input Input) (Input, error) {
	fields := ValidateInput(input)
	for _, key := range []string{"title", "description", "location", "starts_at", "ends_at"} {
		if message, ok := fields[key]; ok {
			return Input{}, invalid(key, message)
		}
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Location = strings.TrimSpace(input.Location)
	input.StartsAt = input.StartsAt.UTC().Truncate(time.Second)
	input.EndsAt = input.EndsAt.UTC().Truncate(time.Second)
	return input, nil
}
func invalid(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrInvalidInput, field, message)
}
