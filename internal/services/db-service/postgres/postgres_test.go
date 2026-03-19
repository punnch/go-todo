package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/punnch/go-todo/internal/core/domains/todo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDB(t *testing.T) (*PostgresRepo, func()) {
	ctx := context.Background()

	// req instance
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-bookworm",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	// get container credentials
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	// create dsn
	dsn := fmt.Sprintf("postgres://test:test@%v:%v/testdb", host, port.Port())

	// create pool
	pool, err := NewPostrgresPool(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	// run migration
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(50) UNIQUE NOT NULL,
			description VARCHAR(200) NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}

	// create repo
	repo := NewPostgresRepo(pool)

	cleanup := func() {
		pool.Close()
		container.Terminate(ctx)
	}

	return repo, cleanup
}

func TestRepo_Create(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	task, err := repo.Create(context.Background(), todo.Task{
		Title:       "programming",
		Description: "5 hours",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID == 0 {
		t.Errorf("expected id to be set, got 0")
	}

	if task.Title != "programming" {
		t.Errorf("exepcted title 'programming', got %v", err)
	}

	if task.Completed {
		t.Errorf("expected completed to be false")
	}
}

func TestRepo_CreateDuplicate(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	task := todo.Task{Title: "programming", Description: "5 hours"}

	_, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.Create(ctx, task)
	if err == nil {
		t.Errorf("expected error on duplicate title, got nil")
	}
}

func TestRepo_GetAll(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	repo.Create(ctx, todo.Task{Title: "task1", Description: "task1"})
	repo.Create(ctx, todo.Task{Title: "task2", Description: "task2"})

	tasks, err := repo.GetAll(ctx, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestRepo_GetAll_FilterByID(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	task1, _ := repo.Create(ctx, todo.Task{Title: "task1", Description: "task1"})
	repo.Create(ctx, todo.Task{Title: "task2", Description: "task2"})

	tasks, err := repo.GetAll(ctx, &task1.ID, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestRepo_GetAll_FilterByCompleted(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()
	completed := true

	task1, _ := repo.Create(ctx, todo.Task{Title: "task1", Description: "task1"})
	repo.Create(ctx, todo.Task{Title: "task2", Description: "task2"})

	repo.Complete(ctx, task1.ID, &completed)

	tasks, err := repo.GetAll(ctx, nil, &completed)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestRepo_Get(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	created, _ := repo.Create(ctx, todo.Task{Title: "task", Description: "title"})

	task, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != created.ID {
		t.Errorf("expected id %d, got %d", created.ID, task.ID)
	}
}

func TestRepo_Get_NotFound(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	_, err := repo.Get(context.Background(), 999)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRepo_Delete(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	created, _ := repo.Create(ctx, todo.Task{Title: "task", Description: "title"})

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.Get(ctx, created.ID)
	if err != pgx.ErrNoRows {
		t.Errorf("expected task to be deleted, got %v", err)
	}
}

func TestRepo_Delete_NotFound(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	err := repo.Delete(context.Background(), 999)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRepo_Complete(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()
	completed := true

	created, _ := repo.Create(ctx, todo.Task{Title: "task", Description: "title"})

	task, err := repo.Complete(ctx, created.ID, &completed)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !task.Completed {
		t.Errorf("expected completed to be true")
	}
}

func TestRepo_Complete_NilCompleted(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	ctx := context.Background()

	created, _ := repo.Create(ctx, todo.Task{Title: "task", Description: "title"})

	task, err := repo.Complete(ctx, created.ID, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Completed != created.Completed {
		t.Errorf("expected completed to remain false")
	}
}

func TestRepo_Complete_NotFound(t *testing.T) {
	repo, cleanup := setupDB(t)
	defer cleanup()

	_, err := repo.Complete(context.Background(), 999, nil)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
