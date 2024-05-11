package services

import (
	"eventom-backend/models"
	"eventom-backend/repositories"
	"eventom-backend/utils"
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
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return us.usersRepository.QuerySignupUser(user.Email, hashedPassword)
}

func (us UsersService) GetUser(email string) (*models.User, *models.ResponseError) {
	return us.usersRepository.QueryGetUser(email)
}

func (us UsersService) ValidateCredentials(user *models.User) (bool, *models.ResponseError) {
	userInDb, responseErr := us.usersRepository.QueryGetUser(user.Email)

	if responseErr != nil {
		return false, responseErr
	}

	if userInDb == nil {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(userInDb.Password), []byte(user.Password))

	if err != nil {
		return false, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		}
	}

	user.ID = userInDb.ID

	return true, nil
}
