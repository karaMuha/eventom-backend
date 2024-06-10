package models

type User struct {
	ID       string `json:"id" validate:"omitempty,uuid"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required"`
}
