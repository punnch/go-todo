package main

import (
	"os"

	core_logger "github.com/punnch/go-todo/api-service/internal/core/logger"
	"github.com/punnch/go-todo/api-service/internal/features/tasks/api"
	"github.com/punnch/go-todo/api-service/internal/features/tasks/client"

	"go.uber.org/zap"
)

func main() {
	// Get log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		panic("LOG_LEVEL not set")
	}

	// Logger
	log, logFileClose, err := core_logger.NewLogger(logLevel, "out/logs") // debug, info, warn, error
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logFileClose()

	// Api service
	apiServiceAddr := os.Getenv("API_SERVICE_ADDR")
	if apiServiceAddr == "" {
		log.Fatal("API_SERVICE_ADDR not set")
	}

	// Service connection
	dbServiceAddr := os.Getenv("DB_SERVICE_ADDR")
	if dbServiceAddr == "" {
		log.Fatal("DB_SERVICE_ADDR not set")
	}

	dbClient := client.NewTodoClient(dbServiceAddr)
	handler := api.NewHandler(dbClient, log)
	srv := api.NewServer(apiServiceAddr, handler)

	if err := srv.StartServer(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
