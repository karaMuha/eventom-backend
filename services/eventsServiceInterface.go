package services

import "eventom-backend/models"

type EventsServiceInterface interface {
	CreateEvent(*models.Event) (*models.Event, *models.ResponseError)

	GetEvent(string) (*models.Event, *models.ResponseError)

	GetAllEvents() ([]*models.Event, *models.ResponseError)

	UpdateEvent(*models.Event) *models.ResponseError
}
