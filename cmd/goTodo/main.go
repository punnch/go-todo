package main

import (
	"context"
	"os"

	"github.com/punnch/go-todo/internal/api/server"
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
	defer pool.Close()

	repo := db.NewPostgresRepo(pool)
	service := todo.NewTodoService(repo)
	handler := server.NewHandler(service)
	router := server.NewRouter(handler)

	if err := server.StartServer(router); err != nil {
		panic(err)
	}
}
