package repositories

import "eventom-backend/models"

type EventsRepositoryInterface interface {
	QueryCreateEvent(event *models.Event) (*models.Event, *models.ResponseError)

	QueryGetEvent(eventId string) (*models.Event, *models.ResponseError)

	QueryGetAllEvents(eventName string, eventLocation string, freeCapacity int) ([]*models.Event, *models.ResponseError)

	QueryUpdateEvent(event *models.Event) (*models.Event, *models.ResponseError)

	QueryIncrementAmountRegistrations(eventId string) (*models.Event, *models.ResponseError)

	QueryDeleteEvent(eventId string) *models.ResponseError
}
