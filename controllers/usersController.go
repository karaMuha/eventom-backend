package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"eventom-backend/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type UsersController struct {
	usersService services.UsersServiceInterface
	validator    *validator.Validate
	logger       *utils.Logger
}

func NewUsersController(usersService services.UsersServiceInterface, logger *utils.Logger) *UsersController {
	return &UsersController{
		usersService: usersService,
		validator:    validator.New(),
		logger:       logger,
	}
}

func (uc UsersController) HandleSignupUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	bodyDecoder := json.NewDecoder(r.Body)

	responseErr := uc.parseUser(&user, bodyDecoder)

	if responseErr != nil {
		uc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseErr = uc.usersService.SignupUser(&user)

	if responseErr != nil {
		uc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	uc.logger.Log(utils.LevelInfo, fmt.Sprintf("User with email %s signed up", user.Email), nil)

	w.WriteHeader(http.StatusOK)
}

func (uc UsersController) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	bodyDecorder := json.NewDecoder(r.Body)

	responseErr := uc.parseUser(&user, bodyDecorder)

	if responseErr != nil {
		uc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	jwtToken, responseErr := uc.usersService.LoginUser(&user)

	if responseErr != nil {
		uc.logger.Log(utils.LevelError, responseErr.Message, nil)
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	uc.logger.Log(utils.LevelInfo, fmt.Sprintf("User with ID %s logged in", user.ID), nil)

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
