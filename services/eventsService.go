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

func (es EventsService) UpdateEvent(userId string, event *models.Event) *models.ResponseError {
	existingEvent, responseErr := es.eventsRepository.QueryGetEvent(event.ID)

	if responseErr != nil {
		return responseErr
	}

	if existingEvent == nil {
		return &models.ResponseError{
			Message: "Event not found",
			Status:  http.StatusNotFound,
		}
	}

	if existingEvent.UserId != userId {
		return &models.ResponseError{
			Message: "Access denied",
			Status:  http.StatusUnauthorized,
		}
	}

	return es.eventsRepository.QueryUpdateEvent(event)
}

func (es EventsService) DeleteEvent(userId string, eventId string) *models.ResponseError {
	event, responseErr := es.eventsRepository.QueryGetEvent(eventId)

	if responseErr != nil {
		return responseErr
	}

	if event == nil {
		return &models.ResponseError{
			Message: "Event not found",
			Status:  http.StatusNotFound,
		}
	}

	if event.UserId != userId {
		return &models.ResponseError{
			Message: "Access denied",
			Status:  http.StatusUnauthorized,
		}
	}

	return es.eventsRepository.QueryDeleteEvent(event.ID)
}
