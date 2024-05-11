package controllers

import (
	"encoding/json"
	"eventom-backend/models"
	"eventom-backend/services"
	"net/http"

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
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uc.validator.Struct(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
