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

func NewRegistrationsController(registrationsService services.RegistrationsServiceInterface, logger *utils.Logger) *RegistrationsController {
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
	userId, ok := r.Context().Value(utils.ContextUserIdKey).(string)

	if !ok {
		http.Error(w, "Could not convert user id from token to a string", http.StatusInternalServerError)
		return
	}

	_, responseErr := rc.registrationsService.CancelRegistration(eventId, userId)

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
