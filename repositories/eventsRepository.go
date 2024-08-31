package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"log"
	"net/http"
	"time"
)

type EventsRepository struct {
	db DBTX
}

func NewEventsRepository(db DBTX) EventsRepositoryInterface {
	return &EventsRepository{
		db: db,
	}
}

func (er *EventsRepository) QueryCreateEvent(event *models.Event) (*models.Event, *models.ResponseError) {
	query := `
		INSERT INTO
			events(event_name, event_description, event_location, event_date, max_capacity, user_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id`
	row := er.db.QueryRow(query, event.Name, event.Description, event.Location, event.Date, event.MaxCapacity, event.UserId)

	var eventId string
	err := row.Scan(&eventId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Event{
		ID:          eventId,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Date:        event.Date,
		UserId:      event.UserId,
	}, nil
}

func (er *EventsRepository) QueryGetEvent(eventId string) (*models.Event, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			events
		WHERE
			id = $1`
	row := er.db.QueryRow(query, eventId)

	var event models.Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.Date, &event.MaxCapacity, &event.AmountRegistration, &event.UserId)

	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &event, nil
}

func (er *EventsRepository) QueryGetAllEvents() ([]*models.Event, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			events`
	rows, err := er.db.Query(query)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	eventsList := make([]*models.Event, 0)
	var eventId, name, description, location, userId string
	var maxCapacity, amountRegistrations int
	var date time.Time

	for rows.Next() {
		err = rows.Scan(&eventId, &name, &description, &location, &date, &maxCapacity, &amountRegistrations, &userId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		event := &models.Event{
			ID:                 eventId,
			Name:               name,
			Description:        description,
			Location:           location,
			Date:               date,
			MaxCapacity:        maxCapacity,
			AmountRegistration: amountRegistrations,
			UserId:             userId,
		}
		eventsList = append(eventsList, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return eventsList, nil
}

func (er *EventsRepository) QueryUpdateEvent(event *models.Event) *models.ResponseError {
	query := `
		UPDATE
			events
		SET
			event_name = $1,
			event_description = $2,
			event_location = $3,
			event_date = $4
		WHERE
			id = $5`
	_, err := er.db.Exec(query, event.Name, event.Description, event.Location, event.Date, event.ID)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (er *EventsRepository) QueryIncrementAmountRegistrations(eventId string) (*models.Event, *models.ResponseError) {
	query := `
		UPDATE
			events
		SET
			amount_registrations = amount_registrations + 1
		WHERE
			id = $1
		RETURNING *`
	row := er.db.QueryRow(query, eventId)

	var event models.Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.Date, &event.MaxCapacity, &event.AmountRegistration, &event.UserId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &event, nil

}

func (er *EventsRepository) QueryDeleteEvent(eventId string) *models.ResponseError {
	query := `
		DELETE FROM
			events
		WHERE
			id = $1`
	_, err := er.db.Exec(query, eventId)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
