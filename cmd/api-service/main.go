package main

import (
	"os"

	"github.com/punnch/go-todo/internal/core/logger"
	"github.com/punnch/go-todo/internal/services/api-service/api/server"
	"github.com/punnch/go-todo/internal/services/api-service/client"

	"go.uber.org/zap"
)

func main() {
	// Get log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		panic("LOG_LEVEL not set")
	}

	// Logger
	log, logFileClose, err := logger.NewLogger(logLevel, "out/logs/api-service") // debug, info, warn, error
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logFileClose()

	// Api service 
	apiServicePort := os.Getenv("API_SERVICE_PORT")
	if apiServicePort == "" {
		log.Fatal("API_SERVICE_PORT not set")
	}

	// Service connection
	dbServicePort := os.Getenv("DB_SERVICE_PORT")
	if dbServicePort == "" {
		log.Fatal("DB_SERVICE_PORT not set")
	}

	dbClient := client.NewTodoClient(dbServicePort)
	handler := server.NewHandler(dbClient, log)
	srv := server.NewServer(apiServicePort, handler)

	if err := srv.StartServer(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
