package main

import (
	"context"

	"github.com/k0kubun/pp"
	"github.com/punnch/go-todo/internal/db"
)

func main() {
	ctx := context.Background()

	pspool := db.NewPostgresPool(ctx)

	conn, err := pspool.CreateConnection()
	if err != nil {
		panic(err)
	}

	// task := todo.NewTask("random title", "random description")

	// dbTask, err := db.Create(ctx, conn, task)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(dbTask)

	rows, err := pspool.GetAll(conn)
	if err != nil {
		panic(err)
	}

	pp.Println(rows)
}
