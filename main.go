package main

import (
	"eventom-backend/config"
	"eventom-backend/server"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting eventom app")

	log.Println("Reading environment variables")
	appEnvironment := os.Getenv("APP_ENV")

	if appEnvironment == "" {
		log.Fatal("Could not get app environment")
	}

	config := config.ReadEnvFile(appEnvironment)

	log.Println("Initializing database")
	db := server.InitDatabase(config)

	log.Println("Initializinh http server")
	httpServer := server.InitHttpServer(config, db)

	log.Printf("Starting app on port %s", config.GetString("SERVER_PORT"))
	httpServer.Start()
}
