package main

import (
	"context"
	"os"

	"github.com/punnch/go-todo/internal/api/handlers"
	"github.com/punnch/go-todo/internal/db"
	"github.com/punnch/go-todo/internal/todo"
)

func main() {
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

	if err := handlers.StartServer(router); err != nil {
		panic(err)
	}
}
