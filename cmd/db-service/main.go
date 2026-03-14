package main

import (
	"context"
	"fmt"
	"os"

	"github.com/punnch/go-todo/internal/services/db-service/postgres"
	"github.com/punnch/go-todo/internal/services/db-service/server"
	"github.com/punnch/go-todo/internal/services/db-service/todo"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		fmt.Println("DB_URL environment variable isn't declared")

		dsn = os.Getenv("LOCAL_DB_URL")
		if dsn == "" {
			return
		}
	}

	pool, err := postgres.NewPostrgresPool(ctx, dsn)
	if err != nil {
		fmt.Println("pool creation err:", err)
		return
	}

	repo := postgres.NewPostgresRepo(pool)

	service := todo.NewTodoService(repo)

	port := os.Getenv("DB_SERVICE_PORT")
	if port == "" {
		fmt.Println("DB_SERVICE_PORT environment variable isn't declared")

		port = ":8010"
	}
	server := server.NewServer(port, service)

	if err := server.StartSever(); err != nil {
		fmt.Println("server error:", err)
		return
	}
}
