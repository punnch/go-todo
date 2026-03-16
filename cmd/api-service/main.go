package main

import (
	"os"

	"github.com/punnch/go-todo/internal/core/logger"
	"github.com/punnch/go-todo/internal/services/api-service/api/server"
	"github.com/punnch/go-todo/internal/services/api-service/client"

	"go.uber.org/zap"
)

func main() {
	// Logger
	log, logFileClose, err := logger.NewLogger("info", "out/logs/api-service") // debug, info, warn, error
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logFileClose()

	// Service connection
	dbURL := os.Getenv("DB_SERVICE_URL")
	if dbURL == "" {
		url := "http://localhost:8010"
		log.Warn("DB_SERVICE_URL not set", zap.String("url", url))
		dbURL = url
	}

	dbClient := client.NewTodoClient(dbURL)
	handler := server.NewHandler(dbClient, log)
	srv := server.NewServer(":8080", handler)

	if err := srv.StartServer(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
