package controllers

import (
	"encoding/json"
	"eventom-backend/dtos"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"net/http"
	"strconv"
	"strings"

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

	userId, ok := r.Context().Value(utils.ContextUserIdKey).(string)

	if !ok {
		http.Error(w, "Could not convert user id from token to a string", http.StatusInternalServerError)
		return
	}

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
	var eventFilters dtos.EventFilterDto
	eventNameParam := r.URL.Query().Get("name")
	locationParam := r.URL.Query().Get("location")
	freeCapacityParam := r.URL.Query().Get("capacity")
	sortColumnParam := r.URL.Query().Get("column")
	sortOrderParam := r.URL.Query().Get("order")

	if strings.EqualFold(sortColumnParam, "") {
		sortColumnParam = "id"
	}

	if strings.EqualFold(sortOrderParam, "") {
		sortOrderParam = "ASC"
	}

	freeCapacity := 0
	if !strings.EqualFold(freeCapacityParam, "") {
		var err error
		freeCapacity, err = strconv.Atoi(freeCapacityParam)
		if err != nil {
			http.Error(w, "free capacity filter must be empty or a number", http.StatusBadRequest)
			return
		}
	}

	eventFilters.Name = eventNameParam
	eventFilters.Location = locationParam
	eventFilters.FreeCapacity = freeCapacity
	eventFilters.SortColumn = sortColumnParam
	eventFilters.SortOrder = sortOrderParam

	err := ec.validator.Struct(&eventFilters)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventsList, responseErr := ec.eventsService.GetAllEvents(&eventFilters)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(eventsList)

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

	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	updatedEvent, responseErr := ec.eventsService.UpdateEvent(userId, &event)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(updatedEvent)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (ec EventsController) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("id")
	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	responseErr := ec.eventsService.DeleteEvent(userId, eventId)

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
