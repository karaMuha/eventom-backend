package repositories

import "eventom-backend/models"

type RegistrationsRepositoryInterface interface {
	QueryRegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError)

	QueryCancelRegistration(registrationId string) *models.ResponseError
}
