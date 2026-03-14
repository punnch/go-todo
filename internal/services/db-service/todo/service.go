package todo

import (
	"context"
	"errors"

	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/core/domains/todo"

	"github.com/jackc/pgx/v5"
)

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTask(ctx context.Context, title, description string) (todo.Task, error) {
	task := todo.NewTask(title, description)

	dbTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return todo.Task{}, err
	}

	return dbTask, nil
}

func (s *TodoService) GetAllTasks(ctx context.Context, id *int, completed *bool) ([]todo.Task, error) {
	tasks, err := s.repo.GetAll(ctx, id, completed)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TodoService) GetTask(ctx context.Context, id int) (todo.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return todo.Task{}, apperrors.ErrNotFound
		}
		return todo.Task{}, err
	}

	return task, nil
}

func (s *TodoService) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	return err
}

func (s *TodoService) CompleteTask(ctx context.Context, id int, completed *bool) (todo.Task, error) {
	task, err := s.repo.Complete(ctx, id, completed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return todo.Task{}, apperrors.ErrNotFound
		}
		return todo.Task{}, apperrors.ErrNotFound
	}

	return task, nil
}
