package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"log"
	"net/http"
)

type RegistrationsRepository struct {
	db DBTX
}

func NewRegistrationsRepository(db DBTX) RegistrationsRepositoryInterface {
	return &RegistrationsRepository{
		db: db,
	}
}

func (rr *RegistrationsRepository) QueryGetRegistration(eventId string, userId string) (*models.Registration, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			registrations
		WHERE
			event_id = $1
			AND
			user_id = $2`
	row := rr.db.QueryRow(query, eventId, userId)

	var registration models.Registration
	err := row.Scan(&registration.ID, &registration.EventId, &registration.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &registration, nil
}

func (rr *RegistrationsRepository) QueryGetAllRegistrations() ([]*models.Registration, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			registrations`
	rows, err := rr.db.Query(query)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	registrationsList := make([]*models.Registration, 0)
	var registrationId, eventId, userId string

	for rows.Next() {
		err = rows.Scan(&registrationId, &eventId, &userId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		registration := &models.Registration{
			ID:      registrationId,
			EventId: eventId,
			UserId:  userId,
		}
		registrationsList = append(registrationsList, registration)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return registrationsList, nil
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
		log.Println(err.Error())
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

func (rr *RegistrationsRepository) QueryCancelRegistration(registrationId string) (*models.Registration, *models.ResponseError) {
	query := `
		DELETE FROM
			registrations
		WHERE
			id := $1
		RETURNING
			*`
	row := rr.db.QueryRow(query, registrationId)
	var deletedRegistration models.Registration

	err := row.Scan(
		&deletedRegistration.ID,
		&deletedRegistration.EventId,
		&deletedRegistration.UserId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ResponseError{
				Message: "Registration not found",
				Status:  http.StatusNotFound,
			}
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &deletedRegistration, nil
}

var _ RegistrationsRepositoryInterface = (*RegistrationsRepository)(nil)
