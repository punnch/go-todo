package main

import (
	"fmt"
	"os"

	"github.com/punnch/go-todo/internal/services/api-service/api/server"
	"github.com/punnch/go-todo/internal/services/api-service/client"
)

func main() {
	dbURL := os.Getenv("DB_SERVICE_URL")
	if dbURL == "" {
		fmt.Println("DB_SERVICE_URL environment variable isn't declared")

		dbURL = "http://localhost:8010"
	}

	dbClient := client.NewTodoClient(dbURL)

	handler := server.NewHandler(dbClient)

	server := server.NewServer(":8080", handler)

	if err := server.Start(); err != nil {
		fmt.Println("server error:", err)
		return
	}
}
