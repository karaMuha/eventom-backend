package models

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"status_code"`
}
