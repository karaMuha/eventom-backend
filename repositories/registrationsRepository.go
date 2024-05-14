package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"net/http"
)

type RegistrationsRepository struct {
	db *sql.DB
}

func NewRegistrationsRepository(db *sql.DB) RegistrationsRepositoryInterface {
	return &RegistrationsRepository{
		db: db,
	}
}

func (rr *RegistrationsRepository) QueryRegisterUserForEvent(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	query := `
		INSERT INTO
			registrations(event_id, user_id)
		VALUES
			($1, $2)
		RETURNING
			id`
	row := rr.db.QueryRow(query, eventId, userId)

	var registrationId string
	err := row.Scan(&registrationId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Registration{
		ID:      registrationId,
		EventId: eventId,
		UserId:  userId,
	}, nil
}

func (rr *RegistrationsRepository) QueryCancelRegistration(registrationId string) *models.ResponseError {
	query := `
		DELETE FROM
			registrations
		WHERE
			id := $1`
	_, err := rr.db.Exec(query, registrationId)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
