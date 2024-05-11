package repositories

import "eventom-backend/models"

type UsersRepositoryInterface interface {
	QuerySignupUser(string, string) *models.ResponseError

	QueryGetUser(string) (*models.User, *models.ResponseError)
}
