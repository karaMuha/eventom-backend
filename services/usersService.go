package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewUsersService(usersRepository repositories.UsersRepositoryInterface) UsersServiceInterface {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) SignupUser(user *models.User) *models.ResponseError {
	hashedPassword, responseErr := hashPassword(user.Password)

	if responseErr != nil {
		return responseErr
	}

	return us.usersRepository.QuerySignupUser(user.Email, hashedPassword)
}

func (us UsersService) GetUser(email string) (*models.User, *models.ResponseError) {
	return us.usersRepository.QueryGetUser(email)
}

func hashPassword(password string) (string, *models.ResponseError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return string(hashedPassword), nil
}
