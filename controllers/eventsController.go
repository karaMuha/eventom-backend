package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type EventsController struct {
	eventsService services.EventsServiceInterface
	validator     *validator.Validate
}

func NewEventsController(eventsService services.EventsServiceInterface) *EventsController {
	return &EventsController{
		eventsService: eventsService,
		validator:     validator.New(),
	}
}

func (ec EventsController) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(utils.ContextUserIdKey).(string)
	event.UserId = userId
	err = ec.validator.Struct(&event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEvent, responseErr := ec.eventsService.CreateEvent(&event)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&createdEvent)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (ec EventsController) HandleGetEvent(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("id")

	event, responseErr := ec.eventsService.GetEvent(eventId)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	responseJson, err := json.Marshal(&event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (ec EventsController) HandleGetAllEvents(w http.ResponseWriter, r *http.Request) {
	eventsList, responseErr := ec.eventsService.GetAllEvents()

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&eventsList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (ec EventsController) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	bodyDecorder := json.NewDecoder(r.Body)
	eventId := r.PathValue("id")
	event.ID = eventId
	responseErr := ec.parseEvent(&event, bodyDecorder)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	existingEvent, responseErr := ec.eventsService.GetEvent(event.ID)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if existingEvent == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	event.UserId = existingEvent.UserId

	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	if event.UserId != userId {
		http.Error(w, "Not allowed", http.StatusUnauthorized)
		return
	}

	responseErr = ec.eventsService.UpdateEvent(&event)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ec EventsController) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("id")

	event, responseErr := ec.eventsService.GetEvent(eventId)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	if event.UserId != userId {
		http.Error(w, "Not allowed", http.StatusUnauthorized)
		return
	}

	responseErr = ec.eventsService.DeleteEvent(event)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ec EventsController) parseEvent(event *models.Event, bodyDecoder *json.Decoder) *models.ResponseError {
	err := bodyDecoder.Decode(event)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	err = ec.validator.Struct(event)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
