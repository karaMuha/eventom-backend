package server

import (
	"database/sql"
	"eventom-backend/config"
	"eventom-backend/controllers"
	"eventom-backend/middlewares"
	"eventom-backend/repositories"
	"eventom-backend/services"
	"eventom-backend/utils"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type HttpServer struct {
	config *viper.Viper
	server *http.Server
}

func InitHttpServer(viperConfig *viper.Viper, db *sql.DB) HttpServer {
	privateKey, err := utils.ReadPrivateKeyFromFile(viperConfig.GetString("PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatalf("Error while reading private key: %v", err)
	}
	config.PrivateKey = privateKey
	eventsRepository := repositories.NewEventsRepository(db)
	usersRepository := repositories.NewUsersRepository(db)
	eventsService := services.NewEventsService(eventsRepository)
	usersService := services.NewUsersService(usersRepository)
	eventsController := controllers.NewEventsController(eventsService)
	usersController := controllers.NewUsersController(usersService)

	router := http.NewServeMux()
	authRouter := http.NewServeMux()

	authRouter.HandleFunc("POST /events", eventsController.HandleCreateEvent)
	router.HandleFunc("GET /events/{id}", eventsController.HandleGetEvent)
	router.HandleFunc("GET /events", eventsController.HandleGetAllEvents)
	authRouter.HandleFunc("PUT /events/{id}", eventsController.HandleUpdateEvent)
	authRouter.HandleFunc("DELETE /events/{id}", eventsController.HandleDeleteEvent)

	router.HandleFunc("POST /signup", usersController.HandleSignupUser)
	router.HandleFunc("POST /login", usersController.HandleLoginUser)
	router.HandleFunc("POST /logout", usersController.HandleLogoutUser)

	router.Handle("/", middlewares.AuthMiddleware(authRouter))

	server := &http.Server{
		Addr:    viperConfig.GetString("SERVER_PORT"),
		Handler: router,
	}

	return HttpServer{
		config: viperConfig,
		server: server,
	}
}

func (hs HttpServer) Start() {
	err := hs.server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
