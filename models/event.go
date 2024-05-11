package models

import (
	"time"
)

type Event struct {
	ID          string    `json:"id" validate:"omitempty,uuid"`
	Name        string    `json:"name" validate:"required,max=100"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=255"`
	Location    string    `json:"location" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	UserId      string    `json:"userId"`
}
