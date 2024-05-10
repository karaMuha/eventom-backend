package services

import "eventom-backend/models"

type UsersServiceInterface interface {
	SignupUser(*models.User) *models.ResponseError

	GetUser(string) (*models.User, *models.ResponseError)
}
