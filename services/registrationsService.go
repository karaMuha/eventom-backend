package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
)

type RegistrationsService struct {
	registrationsRepository repositories.RegistrationsRepositoryInterface
	transactionHandler      repositories.TransactionHandler
}

func NewRegistrationsService(registrationsRepository repositories.RegistrationsRepositoryInterface, transactionHandler repositories.TransactionHandler) *RegistrationsService {
	return &RegistrationsService{
		registrationsRepository: registrationsRepository,
		transactionHandler:      transactionHandler,
	}
}

func (rs RegistrationsService) RegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	return rs.transactionHandler.ExecTx(eventId, userId)
}

func (rs RegistrationsService) GetRegistration(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	return rs.registrationsRepository.QueryGetRegistration(eventId, userId)
}

func (rs RegistrationsService) GetAllRegistration() ([]*models.Registration, *models.ResponseError) {
	return rs.registrationsRepository.QueryGetAllRegistrations()
}

func (rs RegistrationsService) CancelRegistration(registrationId string) *models.ResponseError {
	return rs.registrationsRepository.QueryCancelRegistration(registrationId)
}

var _ RegistrationsServiceInterface = (*RegistrationsService)(nil)
