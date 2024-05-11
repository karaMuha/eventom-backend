package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"net/http"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepositoryInterface {
	return &UsersRepository{
		db: db,
	}
}

func (ur UsersRepository) QuerySignupUser(email string, hashedPassword string) *models.ResponseError {
	query := `
		INSERT INTO
			users(email, password)
		VALUES
			($1, $2)`
	_, err := ur.db.Exec(query, email, hashedPassword)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepository) QueryGetUser(email string) (*models.User, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			users
		WHERE
			email = $1`
	row := ur.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &user, nil
}
