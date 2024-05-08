package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
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

	err = ec.validator.Struct(&event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, responseErr := ec.eventsService.CreateEvent(&event)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
