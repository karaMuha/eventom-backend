package services

import "eventom-backend/models"

type EventsServiceInterface interface {
	CreateEvent(*models.Event) (*models.Event, *models.ResponseError)
}
