package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/punnch/go-todo/internal/core/domains/todo"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pool: pool,
	}
}

// Create pool to work without concurrency problems
func NewPostrgresPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
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
		task.Description,
	).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return todo.Task{}, err
	}

	return dbTask, nil
}

func (p *PostgresRepo) GetAll(ctx context.Context, id *int, completed *bool) ([]todo.Task, error) {
	sql := `
	SELECT id, title, description, completed, created_at 
	FROM tasks
	WHERE ($1::int IS NULL OR id=$1) 
	  AND ($2::boolean IS NULL OR completed=$2)
	ORDER BY id DESC;
	`

	rows, err := p.pool.Query(ctx, sql, id, completed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []todo.Task
	for rows.Next() {
		var task todo.Task
		if err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
		); err != nil {
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
	SELECT id, title, description, completed, created_at FROM TASKS
	WHERE id=$1;
	`

	var dbTask todo.Task
	if err := p.pool.QueryRow(ctx, sql, id).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
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

func (p *PostgresRepo) Complete(ctx context.Context, id int, completed *bool) (todo.Task, error) {
	if completed == nil {
		return p.Get(ctx, id)
	}

	sql := `
	UPDATE tasks
	SET completed=$1
	WHERE id=$2
	RETURNING id, title, description, completed, created_at;
	`

	var dbTask todo.Task
	if err := p.pool.QueryRow(ctx, sql, completed, id).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return todo.Task{}, err
	}

	return dbTask, nil
}
