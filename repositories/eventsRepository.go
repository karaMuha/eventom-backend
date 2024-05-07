package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"net/http"
)

type EventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) EventsRepositoryInterface {
	return &EventsRepository{
		db: db,
	}
}

func (er EventsRepository) QueryCreateEvent(event *models.Event) (*models.Event, *models.ResponseError) {
	query := `
		INSERT INTO
			events(event_name, event_description, event_location, event_date)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id`
	row := er.db.QueryRow(query, event.Name, event.Description, event.Location, event.Date)

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
	}, nil
}
