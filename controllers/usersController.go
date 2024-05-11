package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UsersController struct {
	usersService services.UsersServiceInterface
	validator    *validator.Validate
	privateKey   *rsa.PrivateKey
}

func NewUsersController(usersService services.UsersServiceInterface, privateKey *rsa.PrivateKey) *UsersController {
	return &UsersController{
		usersService: usersService,
		validator:    validator.New(),
		privateKey:   privateKey,
	}
}

func (uc UsersController) HandleSignupUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	bodyDecoder := json.NewDecoder(r.Body)

	responseErr := uc.parseUser(&user, bodyDecoder)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	existingUser, responseErr := uc.usersService.GetUser(user.Email)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if existingUser != nil {
		http.Error(w, "Email address already exists", http.StatusConflict)
		return
	}

	responseErr = uc.usersService.SignupUser(&user)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc UsersController) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	bodyDecorder := json.NewDecoder(r.Body)

	responseErr := uc.parseUser(&user, bodyDecorder)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	validCredentials, responseErr := uc.usersService.ValidateCredentials(&user)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if !validCredentials {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	jwtToken, err := utils.GenerateJwt(user.ID, uc.privateKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: jwtToken,
	})

	w.WriteHeader(http.StatusOK)
}

func (uc UsersController) parseUser(user *models.User, bodyDecoder *json.Decoder) *models.ResponseError {
	err := bodyDecoder.Decode(user)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	err = uc.validator.Struct(user)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
