package repositories

import "eventom-backend/models"

type EventsRepositoryInterface interface {
	QueryCreateEvent(event *models.Event) (*models.Event, *models.ResponseError)

	QueryGetEvent(eventId string) (*models.Event, *models.ResponseError)

	QueryGetAllEvents() ([]*models.Event, *models.ResponseError)

	QueryUpdateEvent(event *models.Event) *models.ResponseError

	QueryDeleteEvent(eventId string) *models.ResponseError
}
