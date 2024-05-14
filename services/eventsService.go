package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
)

type EventsService struct {
	eventsRepository repositories.EventsRepositoryInterface
}

func NewEventsService(eventsRepository repositories.EventsRepositoryInterface) EventsServiceInterface {
	return &EventsService{
		eventsRepository: eventsRepository,
	}
}

func (es EventsService) CreateEvent(event *models.Event) (*models.Event, *models.ResponseError) {
	return es.eventsRepository.QueryCreateEvent(event)
}

func (es EventsService) GetEvent(eventId string) (*models.Event, *models.ResponseError) {
	return es.eventsRepository.QueryGetEvent(eventId)
}

func (es EventsService) GetAllEvents() ([]*models.Event, *models.ResponseError) {
	return es.eventsRepository.QueryGetAllEvents()
}

func (es EventsService) UpdateEvent(event *models.Event) *models.ResponseError {
	return es.eventsRepository.QueryUpdateEvent(event)
}

func (es EventsService) DeleteEvent(event *models.Event) *models.ResponseError {
	return es.eventsRepository.QueryDeleteEvent(event.ID)
}
