package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
)

type RegistrationsService struct {
	registrationsRepository repositories.RegistrationsRepositoryInterface
}

func NewRegistrationsService(registrationsRepository repositories.RegistrationsRepositoryInterface) RegistrationsServiceInterface {
	return &RegistrationsService{
		registrationsRepository: registrationsRepository,
	}
}

func (rs RegistrationsService) RegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	return rs.registrationsRepository.QueryRegisterUserForEvent(eventId, userId)
}

func (rs RegistrationsService) GetRegistration(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	return rs.registrationsRepository.QueryGetRegistration(eventId, userId)
}

func (rs RegistrationsService) CancelRegistration(registrationId string) *models.ResponseError {
	return rs.registrationsRepository.QueryCancelRegistration(registrationId)
}
