package server

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func initDatabase(dbDriver string, dbConnection string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbConnection)
	if err != nil {
		log.Printf("Error while initializing database %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error while validation database connection: %v", err)
		return nil, err
	}

	return db, nil
}

func ConnectToDb() *sql.DB {
	var count int
	dbConnection := os.Getenv("DBCONNECTION")
	dbDriver := os.Getenv("DB_DRIVER")

	for {
		dbHandler, err := initDatabase(dbDriver, dbConnection)

		if err == nil {
			return dbHandler
		}

		log.Println("Postgres container not yet ready...")
		count++
		log.Println(count)

		if count > 10 {
			log.Fatalf("Error while initializing database %v", err)
			return nil
		}

		log.Println("Backing off for five seconds...")
		time.Sleep(5 * time.Second)
	}
}
