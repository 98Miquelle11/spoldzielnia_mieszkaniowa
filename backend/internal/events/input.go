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

func validateID(id int64) error {
	if id <= 0 {
		return invalid("id", "musi być dodatnią liczbą całkowitą")
	}
	return nil
}

func normalize(input Input) (Input, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Location = strings.TrimSpace(input.Location)

	switch {
	case input.Title == "":
		return Input{}, invalid("title", "jest wymagany")
	case utf8.RuneCountInString(input.Title) > MaxTitleLength:
		return Input{}, invalid("title", "maksymalnie 200 znaków")
	}

	switch {
	case input.Description == "":
		return Input{}, invalid("description", "jest wymagany")
	case utf8.RuneCountInString(input.Description) > MaxDescriptionLength:
		return Input{}, invalid("description", "maksymalnie 5000 znaków")
	}

	switch {
	case input.Location == "":
		return Input{}, invalid("location", "jest wymagane")
	case utf8.RuneCountInString(input.Location) > MaxLocationLength:
		return Input{}, invalid("location", "maksymalnie 200 znaków")
	}

	if input.StartsAt.IsZero() {
		return Input{}, invalid("starts_at", "jest wymagane")
	}
	if input.EndsAt.IsZero() {
		return Input{}, invalid("ends_at", "jest wymagane")
	}

	input.StartsAt = input.StartsAt.UTC()
	input.EndsAt = input.EndsAt.UTC()

	if input.EndsAt.Before(input.StartsAt) {
		return Input{}, invalid("ends_at", "nie może być wcześniejsze niż starts_at")
	}
	
	input.StartsAt = input.StartsAt.Truncate(time.Second)
	input.EndsAt = input.EndsAt.Truncate(time.Second)

	return input, nil
}

func invalid(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrInvalidInput, field, message)
}
