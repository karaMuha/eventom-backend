package repositories

import "eventom-backend/models"

type RegistrationsRepositoryInterface interface {
	QueryRegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError)

	QueryGetRegistration(eventId string, userId string) (*models.Registration, *models.ResponseError)

	QueryGetAllRegistrations() ([]*models.Registration, *models.ResponseError)

	QueryCancelRegistration(registrationId string) (*models.Registration, *models.ResponseError)
}
