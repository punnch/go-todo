package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/punnch/go-todo/db-service/internal/core/logger"
	"github.com/punnch/go-todo/db-service/internal/features/tasks/repository/postgres"
	"github.com/punnch/go-todo/db-service/internal/features/tasks/service"
	"github.com/punnch/go-todo/db-service/internal/features/tasks/transport"
)

func main() {
	ctx := context.Background()

	// Get log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		panic("LOG_LEVEL not set")
	}

	// Create logger
	log, logFileClose, err := logger.NewLogger(logLevel, "out/logs") // debug, info, warn, error
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

	dsn := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s",
		postgresUser,
		postgresPassword,
		postgresDB,
	)

	// Create db pool
	pool, err := postgres.NewPostrgresPool(ctx, dsn)
	if err != nil {
		log.Fatal("failed to create database pool", zap.Error(err))
	}
	defer pool.Close()

	// Get server port
	addr := os.Getenv("DB_SERVICE_ADDR")
	if addr == "" {
		log.Fatal("DB_SERVICE_ADDR not set")
	}

	repo := postgres.NewPostgresRepo(pool)
	service := service.NewTodoService(repo)
	srv := transport.NewServer(addr, service, log)

	// Start server
	if err := srv.StartSever(); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
