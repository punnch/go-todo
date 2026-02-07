package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/punnch/go-todo/internal/todo"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pool: pool,
	}
}

func (p *PostgresRepo) Create(ctx context.Context, task todo.Task) (todo.Task, error) {
	sql := `
	INSERT INTO tasks (title, description)
	VALUES($1, $2)
	RETURNING id, title, description, completed, created_at;
	`

	var dbTask todo.Task
	if err := p.pool.QueryRow(
		ctx,
		sql,
		task.Title,
		task.Descripton,
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

func (p *PostgresRepo) GetAll(ctx context.Context) ([]todo.Task, error) {
	sql := `
	SELECT id, title, description, completed, created_at FROM tasks
	ORDER BY id DESC;
	`

	rows, err := p.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (p *PostgresRepo) Get(ctx context.Context, id int) (todo.Task, error) {
	sql := `
	SElECT id, title, description, completed, created_at FROM TASKS
	WHERE id=$1;
	`

	var dbTask todo.Task
	if err := p.pool.QueryRow(ctx, sql, id).Scan(
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

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	sql := `
	DELETE FROM tasks
	WHERE id=$1
	RETURNING id;
	`

	var checkId int
	err := p.pool.QueryRow(ctx, sql, id).Scan(&checkId)

	return err
}

func (p *PostgresRepo) Complete(ctx context.Context, id int) (todo.Task, error) {
	sql := `
	UPDATE tasks
	SET completed=TRUE
	WHERE id=$1
	RETURNING id, title, description, completed, created_at;
	`

	var dbTask todo.Task
	if err := p.pool.QueryRow(ctx, sql, id).Scan(
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
