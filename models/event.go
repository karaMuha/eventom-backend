package models

import (
	"time"
)

type Event struct {
	ID                 string    `json:"id" validate:"omitempty,uuid"`
	Name               string    `json:"name" validate:"required,max=100"`
	Description        string    `json:"description,omitempty" validate:"omitempty,max=255"`
	Location           string    `json:"location" validate:"required"`
	Date               time.Time `json:"date" validate:"required"`
	MaxCapacity        int       `json:"max_capacity" validate:"required,gte=1"`
	AmountRegistration int       `json:"amount_registrations"`
	UserId             string    `json:"userId"`
}
