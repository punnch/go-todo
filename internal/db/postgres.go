package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/punnch/go-todo/internal/todo"
)

// todo:
/*
1. Implement pgxpool.Pool
2. Write ID to Task
3. Return ID in Create
4. Make a struct for postgres
*/

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, "postgres://postgres:password@localhost:5432/go-todo")
}

func CreateTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(50) NOT NULL,
		description VARCHAR(200) NOT NULL,
		completed BOOLEAN NOT NULL,
		created_at TIMESTAMP NOT NULL,

		UNIQUE(title)
	);
	`

	_, err := conn.Exec(ctx, sqlQuery)

	return err
}

func Create(ctx context.Context, conn *pgx.Conn, task todo.Task) error {
	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_at)
	VALUES($1, $2, $3, $4);
	`

	_, err := conn.Exec(ctx, sqlQuery, task.Title, task.Descripton, task.Completed, task.CreatedAt)

	return err
}
