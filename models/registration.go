package models

type Registration struct {
	ID      string `json:"id" validate:"omitempty,uuid"`
	EventId string `json:"event_id" validate:"uuid"`
	UserId  string `json:"user_id" validate:"uuid"`
}
