package services

import "eventom-backend/models"

type UsersServiceInterface interface {
	SignupUser(user *models.User) *models.ResponseError

	GetUser(email string) (*models.User, *models.ResponseError)

	LoginUser(user *models.User) (string, *models.ResponseError)
}
