package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/punnch/go-todo/internal/todo"
)

// todo:
/*
1. Implement pgxpool.Pool
*/

type PostgresPool struct {
	ctx context.Context
}

func NewPostgresPool(ctx context.Context) *PostgresPool {
	return &PostgresPool{
		ctx: ctx,
	}
}

func (p *PostgresPool) CreateConnection() (*pgx.Conn, error) {
	return pgx.Connect(p.ctx, "postgres://postgres:password@localhost:5432/go-todo")
}

func (p *PostgresPool) Create(conn *pgx.Conn, task todo.Task) (todo.Task, error) {
	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING id, title, descripiton, completed, created_at;
	`

	var dbTask todo.Task
	if err := conn.QueryRow(
		p.ctx,
		sqlQuery,
		task.Title,
		task.Descripton,
		task.Completed,
		task.CreatedAt,
	).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Descripton,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return todo.Task{}, err
	}

	return dbTask, nil
}

func (p *PostgresPool) GetAll(conn *pgx.Conn) ([]todo.Task, error) {
	sqlQuery := `
	SELECT id, title, description, completed, created_at FROM tasks
	ORDER BY id DESC;
	`

	rows, err := conn.Query(p.ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	// defer (when implement pool)

	var tasks []todo.Task
	for rows.Next() {
		var task todo.Task
		if err = rows.Scan(&task.ID, &task.Title, &task.Descripton, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
