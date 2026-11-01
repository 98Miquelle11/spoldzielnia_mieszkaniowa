package events

import (
	"errors"
	"time"
)

const (
	MaxTitleLength       = 200
	MaxLocationLength    = 200
	MaxDescriptionLength = 5000
)

var (
	ErrNotFound     = errors.New("event not found")
	ErrInvalidInput = errors.New("invalid event input")
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
