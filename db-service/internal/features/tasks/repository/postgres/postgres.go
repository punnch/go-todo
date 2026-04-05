package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/punnch/go-todo/db-service/internal/core/domain"
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

func (p *PostgresRepo) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	sql := `
	INSERT INTO tasks (title, description, completed, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING id, title, description, completed, created_at;
	`

	var dbTask domain.Task
	if err := p.pool.QueryRow(
		ctx,
		sql,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
	).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return domain.Task{}, err
	}

	return dbTask, nil
}

func (p *PostgresRepo) GetAll(ctx context.Context, id *int, completed *bool) ([]domain.Task, error) {
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

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
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

func (p *PostgresRepo) Get(ctx context.Context, id int) (domain.Task, error) {
	sql := `
	SELECT id, title, description, completed, created_at FROM TASKS
	WHERE id=$1;
	`

	var dbTask domain.Task
	if err := p.pool.QueryRow(ctx, sql, id).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return domain.Task{}, err
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

func (p *PostgresRepo) Complete(ctx context.Context, id int, completed *bool) (domain.Task, error) {
	if completed == nil {
		return p.Get(ctx, id)
	}

	sql := `
	UPDATE tasks
	SET completed=$1
	WHERE id=$2
	RETURNING id, title, description, completed, created_at;
	`

	var dbTask domain.Task
	if err := p.pool.QueryRow(ctx, sql, completed, id).Scan(
		&dbTask.ID,
		&dbTask.Title,
		&dbTask.Description,
		&dbTask.Completed,
		&dbTask.CreatedAt,
	); err != nil {
		return domain.Task{}, err
	}

	return dbTask, nil
}
