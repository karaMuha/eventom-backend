package models

import "time"

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location"`
	Date        time.Time `json:"date"`
}
