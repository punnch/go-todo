package main

import (
	"context"

	"github.com/punnch/go-todo/internal/db"
	"github.com/punnch/go-todo/internal/todo"
)

func main() {
	ctx := context.Background()

	conn, err := db.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	if err := db.CreateTable(ctx, conn); err != nil {
		panic(err)
	}

	task := todo.NewTask("some", "aa")

	if err := db.Create(ctx, conn, task); err != nil {
		panic(err)
	}
}
