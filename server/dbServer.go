package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	host := config.GetString("DB_HOST")
	port := config.GetString("DB_PORT")
	user := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	dbName := config.GetString("DB_NAME")
	dbDriver := config.GetString("DB_DRIVER")
	sslMode := config.GetString("SSL_MODE")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)

	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		log.Fatalf("Error while initializing database %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while validation database connection: %v", err)
	}

	return db
}
