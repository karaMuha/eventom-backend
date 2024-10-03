package controllers

import (
	"encoding/json"
	"errors"
	"eventom-backend/dtos"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type EventsController struct {
	eventsService services.EventsServiceInterface
	validator     *validator.Validate
	logger        *utils.Logger
}

func NewEventsController(eventsService services.EventsServiceInterface, logger *utils.Logger) *EventsController {
	return &EventsController{
		eventsService: eventsService,
		validator:     validator.New(),
		logger:        logger,
	}
}

func (ec EventsController) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		ec.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(utils.ContextUserIdKey).(string)

	if !ok {
		ec.logger.Log(utils.LevelFatal, "Could not convert user id from token to a string", nil)
		http.Error(w, "Could not convert user id from token to a string", http.StatusInternalServerError)
		return
	}

	event.UserId = userId
	err = ec.validator.Struct(&event)

	if err != nil {
		ec.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEvent, responseErr := ec.eventsService.CreateEvent(&event)

	if responseErr != nil {
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	ec.logger.Log(utils.LevelInfo, fmt.Sprintf("Event with ID %s created", createdEvent.ID), nil)

	responseJson, err := json.Marshal(&createdEvent)

	if err != nil {
		ec.logger.Log(utils.LevelFatal, err.Error(), nil)
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
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&event)

	if err != nil {
		ec.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (ec EventsController) HandleGetAllEvents(w http.ResponseWriter, r *http.Request) {
	eventFilters, err := setEventFilters(r)

	if err != nil {
		ec.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ec.validator.Struct(eventFilters)

	if err != nil {
		ec.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventsList, totalCount, responseErr := ec.eventsService.GetAllEvents(eventFilters)

	if responseErr != nil {
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	eventListMetadata := &dtos.EventListMetadata{
		CurrentPage:  eventFilters.Page,
		PageSize:     eventFilters.PageSize,
		LastPage:     int(math.Ceil(float64(totalCount) / float64(eventFilters.PageSize))),
		TotalRecords: totalCount,
	}

	responseData := &dtos.EventListResponse{
		Events:   eventsList,
		Metadata: eventListMetadata,
	}

	responseJson, err := json.Marshal(responseData)

	if err != nil {
		ec.logger.Log(utils.LevelFatal, err.Error(), nil)
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
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	updatedEvent, responseErr := ec.eventsService.UpdateEvent(userId, &event)

	if responseErr != nil {
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	ec.logger.Log(utils.LevelInfo, fmt.Sprintf("Event with ID %s updated", updatedEvent.ID), nil)

	responseJson, err := json.Marshal(updatedEvent)

	if err != nil {
		ec.logger.Log(utils.LevelFatal, err.Error(), nil)
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
		ec.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	ec.logger.Log(utils.LevelInfo, fmt.Sprintf("Event with ID %s deleted", eventId), nil)

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

func setEventFilters(r *http.Request) (*dtos.EventFilterDto, error) {
	var eventFilters dtos.EventFilterDto
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("page_size")
	eventNameParam := r.URL.Query().Get("name")
	locationParam := r.URL.Query().Get("location")
	freeCapacityParam := r.URL.Query().Get("capacity")
	sortColumnParam := r.URL.Query().Get("column")
	sortOrderParam := r.URL.Query().Get("order")

	page := 1
	if !strings.EqualFold(pageParam, "") {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			return nil, err
		}
		if page < 1 {
			return nil, errors.New("page must be greater or equal 1")
		}
	}

	pageSize := 10
	if !strings.EqualFold(pageSizeParam, "") {
		var err error
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil {
			return nil, errors.New("page size must be a number")
		}
	}

	freeCapacity := 0
	if !strings.EqualFold(freeCapacityParam, "") {
		var err error
		freeCapacity, err = strconv.Atoi(freeCapacityParam)
		if err != nil {
			return nil, errors.New("free capacity must be empty or a number")
		}
	}

	if strings.EqualFold(sortColumnParam, "") {
		sortColumnParam = "id"
	}

	if strings.EqualFold(sortOrderParam, "") {
		sortOrderParam = "ASC"
	}

	eventFilters.Name = eventNameParam
	eventFilters.Location = locationParam
	eventFilters.FreeCapacity = freeCapacity
	eventFilters.SortColumn = sortColumnParam
	eventFilters.SortOrder = sortOrderParam
	eventFilters.Page = page
	eventFilters.PageSize = pageSize

	return &eventFilters, nil
}
