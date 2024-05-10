package repositories

import "eventom-backend/models"

type EventsRepositoryInterface interface {
	QueryCreateEvent(*models.Event) (*models.Event, *models.ResponseError)

	QueryGetEvent(string) (*models.Event, *models.ResponseError)

	QueryGetAllEvents() ([]*models.Event, *models.ResponseError)

	QueryUpdateEvent(*models.Event) *models.ResponseError
}
