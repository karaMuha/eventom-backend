package models

type Registration struct {
	ID      string `json:"id" validate:"omitempty,uuid"`
	EventId string `json:"eventId" validate:"uuid"`
	UserId  string `json:"userId" validate:"uuid"`
}
