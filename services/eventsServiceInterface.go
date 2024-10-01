package services

import (
	"eventom-backend/dtos"
	"eventom-backend/models"
)

type EventsServiceInterface interface {
	CreateEvent(event *models.Event) (*models.Event, *models.ResponseError)

	GetEvent(eventId string) (*models.Event, *models.ResponseError)

	GetAllEvents(eventFilters *dtos.EventFilterDto) ([]*models.Event, *models.ResponseError)

	UpdateEvent(userId string, event *models.Event) (*models.Event, *models.ResponseError)

	DeleteEvent(userId string, eventId string) *models.ResponseError
}
