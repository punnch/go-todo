package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/punnch/go-todo/internal/api/handlers"
	"github.com/punnch/go-todo/internal/db"
	"github.com/punnch/go-todo/internal/todo"
)

// todo:
// 1. env vars
// 2. http routing
// 3. migrations

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DB_URL")

	ctx := context.Background()

	pool, err := db.NewPostrgresPool(ctx, dbURL)
	if err != nil {
		panic(err)
	}

	repo := db.NewPostgresRepo(pool)
	service := todo.NewTodoService(repo)
	handler := handlers.NewHandler(service)
	router := handlers.NewRouter(handler)

	handlers.StartServer(router)
}
