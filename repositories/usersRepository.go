package repositories

import (
	"database/sql"
	"eventom-backend/models"
	"net/http"
	"strings"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (ur *UsersRepository) QuerySignupUser(email string, hashedPassword string) *models.ResponseError {
	query := `
		INSERT INTO
			users(email, password)
		VALUES
			($1, $2)`
	_, err := ur.db.Exec(query, email, hashedPassword)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return &models.ResponseError{
				Message: "Email already registered",
				Status:  http.StatusConflict,
			}
		}
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur *UsersRepository) QueryGetUser(email string) (*models.User, *models.ResponseError) {
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
			return nil, &models.ResponseError{
				Message: "User not found",
				Status:  http.StatusNotFound,
			}
		}
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &user, nil
}

var _ UsersRepositoryInterface = (*UsersRepository)(nil)
