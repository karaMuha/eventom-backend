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

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	registrationsRepository := NewRegistrationsRepository(tx)
	eventsRepository := NewEventsRepository(tx)

	registration, responseErr := registrationsRepository.QueryRegisterUserForEvent(evntId, userId)

	if responseErr != nil {
		tx.Rollback()
		return nil, responseErr
	}



	_ = tx.Commit()

	return registration, nil
}
