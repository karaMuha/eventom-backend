package services

import "eventom-backend/models"

type RegistrationsServiceInterface interface {
	RegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError)

	GetRegistration(eventId string, userId string) (*models.Registration, *models.ResponseError)

	GetAllRegistration() ([]*models.Registration, *models.ResponseError)

	CancelRegistration(registrationId string) *models.ResponseError
}
