package todo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/core/domains/todo"
)

type mockRepo struct {
	createFn   func(ctx context.Context, task todo.Task) (todo.Task, error)
	getAllFn   func(ctx context.Context, id *int, completed *bool) ([]todo.Task, error)
	getFn      func(ctx context.Context, id int) (todo.Task, error)
	deleteFn   func(ctx context.Context, id int) error
	completeFn func(ctx context.Context, id int, completed *bool) (todo.Task, error)
}

func (m *mockRepo) Create(ctx context.Context, task todo.Task) (todo.Task, error) {
	return m.createFn(ctx, task)
}

func (m *mockRepo) GetAll(ctx context.Context, id *int, completed *bool) ([]todo.Task, error) {
	return m.getAllFn(ctx, id, completed)
}

func (m *mockRepo) Get(ctx context.Context, id int) (todo.Task, error) {
	return m.getFn(ctx, id)
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func (m *mockRepo) Complete(ctx context.Context, id int, completed *bool) (todo.Task, error) {
	return m.completeFn(ctx, id, completed)
}

func TestCreateTask(t *testing.T) {
	repo := &mockRepo{
		createFn: func(ctx context.Context, task todo.Task) (todo.Task, error) {
			task.ID = 1
			return task, nil
		},
	}

	service := NewTodoService(repo)
	task, err := service.CreateTask(context.Background(), "programming", "5 hours")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected task id 1, got %d", task.ID)
	}
}

func TestGetAllTasks(t *testing.T) {
	completed := true

	repo := &mockRepo{
		getAllFn: func(ctx context.Context, id *int, completed *bool) ([]todo.Task, error) {
			return []todo.Task{
				{ID: 1, Title: "programming", Description: "5 hours", Completed: true},
				{ID: 5, Title: "wake up", Description: "in 6:30 AM", Completed: true},
			}, nil
		},
	}

	service := NewTodoService(repo)
	tasks, err := service.GetAllTasks(context.Background(), nil, &completed)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetTask(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, id int) (todo.Task, error) {
			return todo.Task{
				ID:          id,
				Title:       "english",
				Description: "read a book",
				Completed:   true,
				CreatedAt:   time.Now(),
			}, nil
		},
	}

	service := NewTodoService(repo)
	task, err := service.GetTask(context.Background(), 3)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != 3 {
		t.Errorf("expected task id 3, got %d", task.ID)
	}

	if task.Title != "english" {
		t.Errorf("expected title 'english', got %s", task.Title)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, id int) (todo.Task, error) {
			return todo.Task{}, pgx.ErrNoRows
		},
	}

	service := NewTodoService(repo)
	_, err := service.GetTask(context.Background(), 999)

	if !errors.Is(err, apperrors.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	service := NewTodoService(repo)
	err := service.DeleteTask(context.Background(), 3)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(ctx context.Context, id int) error {
			return pgx.ErrNoRows
		},
	}

	service := NewTodoService(repo)
	err := service.DeleteTask(context.Background(), 2)

	if !errors.Is(err, apperrors.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestCompleteTask(t *testing.T) {
	completed := true

	repo := &mockRepo{
		completeFn: func(ctx context.Context, id int, completed *bool) (todo.Task, error) {
			return todo.Task{
				ID:        id,
				Completed: *completed,
			}, nil
		},
	}

	service := NewTodoService(repo)
	task, err := service.CompleteTask(context.Background(), 3, &completed)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != 3 {
		t.Errorf("expected task id 3, got %d", task.ID)
	}

	if !task.Completed {
		t.Errorf("expected task to be completed")
	}
}

func TestCompleteTask_NotFound(t *testing.T) {
	repo := &mockRepo{
		completeFn: func(ctx context.Context, id int, completed *bool) (todo.Task, error) {
			return todo.Task{}, pgx.ErrNoRows
		},
	}

	service := NewTodoService(repo)
	_, err := service.CompleteTask(context.Background(), 2, nil)

	if !errors.Is(err, apperrors.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
