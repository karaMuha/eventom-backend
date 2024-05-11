package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
	"net/http"
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

func (es EventsService) UpdateEvent(event *models.Event, userId string) *models.ResponseError {
	if event.UserId != userId {
		return &models.ResponseError{
			Message: "Only the author of an event can edit the event",
			Status:  http.StatusUnauthorized,
		}
	}

	return es.eventsRepository.QueryUpdateEvent(event)
}

func (es EventsService) DeleteEvent(event *models.Event, userId string) *models.ResponseError {
	if event.UserId != userId {
		return &models.ResponseError{
			Message: "Only the author of an event can edit the event",
			Status:  http.StatusUnauthorized,
		}
	}

	return es.eventsRepository.QueryDeleteEvent(event.ID)
}
