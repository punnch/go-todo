package main

import (
	"context"
	"os"

	"github.com/punnch/go-todo/internal/core/logger"
	"github.com/punnch/go-todo/internal/services/db-service/postgres"
	"github.com/punnch/go-todo/internal/services/db-service/server"
	"github.com/punnch/go-todo/internal/services/db-service/todo"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Get log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		panic("LOG_LEVEL not set")
	}

	// Create logger
	log, logFileClose, err := logger.NewLogger(logLevel, "out/logs/db-service") // debug, info, warn, error
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logFileClose()

	// Get db credentials
	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
			log.Fatal("POSTGRES_USER not set")
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
			log.Fatal("POSTGRES_PASSWORD not set")
	}

	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresDB == "" {
			log.Fatal("POSTGRES_DB not set")
	}

	// Create db pool
	pool, err := postgres.NewPostrgresPool(ctx, postgresUser, postgresPassword, postgresDB); if err != nil {
		log.Fatal("failed to create database pool", zap.Error(err))
	}
	defer pool.Close()

	// Get server port
	port := os.Getenv("DB_SERVICE_PORT")
	if port == "" {
		log.Fatal("DB_SERVICE_PORT not set")
	}

	repo := postgres.NewPostgresRepo(pool)
	service := todo.NewTodoService(repo)
	srv := server.NewServer(port, service, log)

	// Start server
	if err := srv.StartSever(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
