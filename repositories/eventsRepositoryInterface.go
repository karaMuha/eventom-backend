package repositories

import "eventom-backend/models"

type EventsRepositoryInterface interface {
	QueryCreateEvent(event *models.Event) (*models.Event, *models.ResponseError)
}
