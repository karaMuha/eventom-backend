package repositories

import "eventom-backend/models"

type UsersRepositoryInterface interface {
	QuerySignupUser(email string, password string) *models.ResponseError

	QueryGetUser(email string) (*models.User, *models.ResponseError)
}
