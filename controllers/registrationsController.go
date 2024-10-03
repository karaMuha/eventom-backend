package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RegistrationsController struct {
	registrationsService services.RegistrationsServiceInterface
	validator            *validator.Validate
	logger               *utils.Logger
}

func NewRegistrationsController(registrationsService services.RegistrationsServiceInterface, logger *utils.Logger) *RegistrationsController {
	return &RegistrationsController{
		registrationsService: registrationsService,
		validator:            validator.New(),
		logger:               logger,
	}
}

func (rc RegistrationsController) HandleRegisterUserForEvent(w http.ResponseWriter, r *http.Request) {
	var registration models.Registration
	err := json.NewDecoder(r.Body).Decode(&registration)

	if err != nil {
		rc.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(utils.ContextUserIdKey).(string)
	if !ok {
		rc.logger.Log(utils.LevelFatal, "Could not convert user id from token to a string", nil)
		http.Error(w, "Could not convert user id from token to a string", http.StatusInternalServerError)
		return
	}

	registration.UserId = userId
	err = rc.validator.Struct(&registration)

	if err != nil {
		rc.logger.Log(utils.LevelError, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdRegistration, responseErr := rc.registrationsService.RegisterUserForEvent(registration.EventId, registration.UserId)

	if responseErr != nil {
		rc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	rc.logger.Log(utils.LevelInfo, fmt.Sprintf("Registration with ID %s created", createdRegistration.ID), nil)

	responseJson, err := json.Marshal(&createdRegistration)

	if err != nil {
		rc.logger.Log(utils.LevelFatal, err.Error(), nil)
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
		rc.logger.Log(utils.LevelFatal, "Could not convert user id from token to a string", nil)
		http.Error(w, "Could not convert user id from token to a string", http.StatusInternalServerError)
		return
	}

	cancelledRegistration, responseErr := rc.registrationsService.CancelRegistration(eventId, userId)

	if responseErr != nil {
		rc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	rc.logger.Log(utils.LevelInfo, fmt.Sprintf("Registration with ID %s cancelled", cancelledRegistration.ID), nil)

	w.WriteHeader(http.StatusOK)
}

func (rc RegistrationsController) HandleGetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	registrationsList, responseErr := rc.registrationsService.GetAllRegistration()

	if responseErr != nil {
		rc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(&registrationsList)

	if err != nil {
		rc.logger.Log(utils.LevelFatal, err.Error(), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
