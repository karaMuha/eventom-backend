package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RegistrationsController struct {
	registrationsService services.RegistrationsServiceInterface
	validator            *validator.Validate
}

func NewRegistrationsController(registrationsService services.RegistrationsServiceInterface) *RegistrationsController {
	return &RegistrationsController{
		registrationsService: registrationsService,
		validator:            validator.New(),
	}
}

func (rc RegistrationsController) HandleRegisterUserForEvent(w http.ResponseWriter, r *http.Request) {
	var registration models.Registration
	err := json.NewDecoder(r.Body).Decode(&registration)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(utils.ContextUserIdKey).(string)
	registration.UserId = userId
	err = rc.validator.Struct(&registration)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingEvent, responseErr := rc.registrationsService.GetRegistration(registration.EventId, userId)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if existingEvent != nil {
		http.Error(w, "User is already registered for this event", http.StatusConflict)
		return
	}

	createdRegistration, responseErr := rc.registrationsService.RegisterUserForEvent(registration.EventId, registration.UserId)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&createdRegistration)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (rc RegistrationsController) HandleCancleRegistration(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("id")
	userId := r.Context().Value(utils.ContextUserIdKey).(string)

	registration, responseErr := rc.registrationsService.GetRegistration(eventId, userId)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if registration == nil {
		http.Error(w, "Registration not found", http.StatusNotFound)
		return
	}

	_, responseErr = rc.registrationsService.CancelRegistration(registration.ID)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc RegistrationsController) HandleGetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	registrationsList, responseErr := rc.registrationsService.GetAllRegistration()

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&registrationsList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
