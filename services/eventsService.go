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
