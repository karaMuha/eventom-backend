package main

import (
	"eventom-backend/server"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting eventom app")

	log.Println("Initializing database")
	db := server.ConnectToDb()

	log.Println("Initializinh http server")
	httpServer := server.InitHttpServer(db)

	log.Printf("Starting app on port %s", os.Getenv("SERVER_PORT"))
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
