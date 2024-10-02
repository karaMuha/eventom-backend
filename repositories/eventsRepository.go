package repositories

import (
	"database/sql"
	"eventom-backend/dtos"
	"eventom-backend/models"
	"fmt"
	"net/http"
	"time"
)

type EventsRepository struct {
	db DBTX
}

func NewEventsRepository(db DBTX) *EventsRepository {
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
		if err == sql.ErrNoRows {
			return nil, &models.ResponseError{
				Message: "Event not found",
				Status:  http.StatusNotFound,
			}
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &event, nil
}

func (er *EventsRepository) QueryGetAllEvents(eventFilters *dtos.EventFilterDto) ([]*models.Event, int, *models.ResponseError) {
	// TODO: checkout squirrel for conditional query building on runtime so the query only has the parts it needs to run. That might improve caching performance
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) OVER(),
			*
		FROM
			events
		WHERE
			(to_tsvector('simple', event_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND
			(event_location = $2 OR $2 = '')
			AND
			((((max_capacity - amount_registrations) >= $3) AND $3 != 0) OR $3 = 0)
		ORDER BY
			%s %s, id ASC
		LIMIT
			$4
		OFFSET
			$5`, eventFilters.SortColumn, eventFilters.SortOrder)
	offset := eventFilters.PageSize * (eventFilters.Page - 1)
	rows, err := er.db.Query(query, eventFilters.Name, eventFilters.Location, eventFilters.FreeCapacity, eventFilters.PageSize, offset)

	if err != nil {
		return nil, 0, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	totalcount := 0
	eventsList := make([]*models.Event, 0)
	var eventId, name, description, location, userId string
	var maxCapacity, amountRegistrations int
	var date time.Time

	for rows.Next() {
		err = rows.Scan(&totalcount, &eventId, &name, &description, &location, &date, &maxCapacity, &amountRegistrations, &userId)
		if err != nil {
			return nil, 0, &models.ResponseError{
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
		return nil, 0, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return eventsList, totalcount, nil
}

func (er *EventsRepository) QueryUpdateEvent(event *models.Event) (*models.Event, *models.ResponseError) {
	query := `
		UPDATE
			events
		SET
			event_name = $1,
			event_description = $2,
			event_location = $3,
			event_date = $4
		WHERE
			id = $5
		RETURNING
			*`
	row := er.db.QueryRow(query, event.Name, event.Description, event.Location, event.Date, event.ID)

	var updatedEvent models.Event
	err := row.Scan(
		&updatedEvent.ID,
		&updatedEvent.Name,
		&updatedEvent.Description,
		&updatedEvent.Location,
		&updatedEvent.Date,
		&updatedEvent.MaxCapacity,
		&updatedEvent.AmountRegistration,
		&updatedEvent.UserId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ResponseError{
				Message: "Event not found",
				Status:  http.StatusNotFound,
			}
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &updatedEvent, nil
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

var _ EventsRepositoryInterface = (*EventsRepository)(nil)
