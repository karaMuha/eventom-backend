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
}

func NewEventsController(eventsService services.EventsServiceInterface) *EventsController {
	return &EventsController{
		eventsService: eventsService,
	}
}

func (ec EventsController) ApiCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = validator.New().Struct(event)

	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
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
