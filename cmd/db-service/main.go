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

	// Create logger
	log, logFileClose, err := logger.NewLogger("info", "out/logs/db-service") // debug, info, warn, error
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logFileClose()

	// Get db url
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Warn("DB_URL not set, trying LOCAL_DB_URL...")

		dsn = os.Getenv("LOCAL_DB_URL")
		if dsn == "" {
			log.Fatal("LOCAL_DB_URL is also not set")
		}
	}

	// Create db pool
	pool, err := postgres.NewPostrgresPool(ctx, dsn)
	if err != nil {
		log.Fatal("failed to create database pool", zap.Error(err))
	}

	// Get server port
	port := os.Getenv("DB_SERVICE_PORT")
	if port == "" {
		port = ":8010"
		log.Warn("DB_SERVICE_PORT not set", zap.String("port", port))
	}

	repo := postgres.NewPostgresRepo(pool)
	service := todo.NewTodoService(repo)
	srv := server.NewServer(port, service, log)

	// Start server
	if err := srv.StartSever(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
