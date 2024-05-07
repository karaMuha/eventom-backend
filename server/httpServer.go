package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type HttpServer struct {
	config *viper.Viper
	server *http.Server
}

func InitHttpServer(config *viper.Viper, db *sql.DB) HttpServer {
	router := http.NewServeMux()

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
