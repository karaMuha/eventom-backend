package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"net/http"
)

type TransactionHandler struct {
	db *sql.DB
}

func NewTxHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{
		db: db,
	}
}

func (th *TransactionHandler) ExecTx(evntId string, userId string) (*models.Registration, *models.ResponseError) {
	tx, err := th.db.Begin()
	registrationsRepository := NewRegistrationsRepository(tx)
	eventsRepository := NewEventsRepository(tx)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	registration, responseErr := registrationsRepository.QueryRegisterUserForEvent(evntId, userId)

	if responseErr != nil {
		tx.Rollback()
		return nil, responseErr
	}

	event, responseErr := eventsRepository.QueryIncrementAmountRegistrations(evntId)

	if responseErr != nil {
		tx.Rollback()
		return nil, responseErr
	}

	if event.AmountRegistration > event.MaxCapacity {
		tx.Rollback()
		return nil, &models.ResponseError{
			Message: "Event is full",
			Status:  http.StatusConflict,
		}
	}

	_ = tx.Commit()

	return registration, nil
}
