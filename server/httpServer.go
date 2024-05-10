package server

import (
	"database/sql"
	"eventom-backend/controllers"
	"eventom-backend/repositories"
	"eventom-backend/services"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type HttpServer struct {
	config *viper.Viper
	server *http.Server
}

func InitHttpServer(config *viper.Viper, db *sql.DB) HttpServer {
	eventsRepository := repositories.NewEventsRepository(db)
	usersRepository := repositories.NewUsersRepository(db)
	eventsService := services.NewEventsService(eventsRepository)
	usersService := services.NewUsersService(usersRepository)
	eventsController := controllers.NewEventsController(eventsService)
	usersController := controllers.NewUsersController(usersService)

	router := http.NewServeMux()

	router.HandleFunc("POST /events", eventsController.HandleCreateEvent)
	router.HandleFunc("GET /events/{id}", eventsController.HandleGetEvent)
	router.HandleFunc("GET /events", eventsController.HandleGetAllEvents)
	router.HandleFunc("PUT /events/{id}", eventsController.HandleUpdateEvent)
	router.HandleFunc("DELETE /events/{id}", eventsController.HandleDeleteEvent)

	router.HandleFunc("POST /signup", usersController.HandleSignupUser)

	server := &http.Server{
		Addr:    config.GetString("SERVER_PORT"),
		Handler: router,
	}

	return HttpServer{
		config: config,
		server: server,
	}
}

func (hs HttpServer) Start() {
	err := hs.server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
