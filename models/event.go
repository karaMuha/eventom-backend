package models

import "time"

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}
