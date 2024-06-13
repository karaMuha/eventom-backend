package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type UsersController struct {
	usersService services.UsersServiceInterface
	validator    *validator.Validate
}

func NewUsersController(usersService services.UsersServiceInterface) *UsersController {
	return &UsersController{
		usersService: usersService,
		validator:    validator.New(),
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

	jwtToken, responseErr := uc.usersService.LoginUser(&user)

	if responseErr != nil {
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}

func (uc UsersController) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now(),
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
