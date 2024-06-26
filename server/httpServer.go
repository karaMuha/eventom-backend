package server

import (
	"database/sql"
	"eventom-backend/controllers"
	"eventom-backend/middlewares"
	"eventom-backend/repositories"
	"eventom-backend/services"
	"eventom-backend/utils"
	"log"
	"net/http"
	"os"
)

func InitHttpServer(db *sql.DB) *http.Server {
	// initialize private key that is used to sign and verify jwt
	err := utils.ReadPrivateKeyFromFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatalf("Error while reading private key: %v", err)
	}

	// initialize protected routes map that is used in auth middleware to determine whether a request needs to be authenticated or not
	utils.SetProtectedRoutes()

	eventsRepository := repositories.NewEventsRepository(db)
	usersRepository := repositories.NewUsersRepository(db)
	registrationsRepository := repositories.NewRegistrationsRepository(db)

	eventsService := services.NewEventsService(eventsRepository)
	usersService := services.NewUsersService(usersRepository)
	registrationsService := services.NewRegistrationsService(registrationsRepository)

	eventsController := controllers.NewEventsController(eventsService)
	usersController := controllers.NewUsersController(usersService)
	registrationsController := controllers.NewRegistrationsController(registrationsService)

	router := http.NewServeMux()

	router.HandleFunc("POST /events", eventsController.HandleCreateEvent)
	router.HandleFunc("GET /events/{id}", eventsController.HandleGetEvent)
	router.HandleFunc("GET /events", eventsController.HandleGetAllEvents)
	router.HandleFunc("PUT /events/{id}", eventsController.HandleUpdateEvent)
	router.HandleFunc("DELETE /events/{id}", eventsController.HandleDeleteEvent)

	router.HandleFunc("POST /signup", usersController.HandleSignupUser)
	router.HandleFunc("POST /login", usersController.HandleLoginUser)
	router.HandleFunc("POST /logout", usersController.HandleLogoutUser)

	router.HandleFunc("POST /events/{id}", registrationsController.HandleRegisterUserForEvent)
	router.HandleFunc("DELETE /registrations/{id}", registrationsController.HandleCancleRegistration)

	return &http.Server{
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: middlewares.AuthMiddleware(router),
	}
}
