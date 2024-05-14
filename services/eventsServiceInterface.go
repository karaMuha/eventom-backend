package services

import "eventom-backend/models"

type EventsServiceInterface interface {
	CreateEvent(event *models.Event) (*models.Event, *models.ResponseError)

	GetEvent(eventId string) (*models.Event, *models.ResponseError)

	GetAllEvents() ([]*models.Event, *models.ResponseError)

	UpdateEvent(event *models.Event) *models.ResponseError

	DeleteEvent(event *models.Event) *models.ResponseError
}
