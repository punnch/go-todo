package service

import (
	"context"
	"errors"

	"github.com/punnch/go-todo/db-service/internal/core/domain"
	core_errors "github.com/punnch/go-todo/db-service/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

type TodoRepository interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	GetAll(ctx context.Context, id *int, completed *bool) ([]domain.Task, error)
	Get(ctx context.Context, id int) (domain.Task, error)
	Delete(ctx context.Context, id int) error
	Complete(ctx context.Context, id int, completed *bool) (domain.Task, error)
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTask(ctx context.Context, title, description string) (domain.Task, error) {
	task := domain.NewTask(title, description)

	dbTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	return dbTask, nil
}

func (s *TodoService) GetAllTasks(ctx context.Context, id *int, completed *bool) ([]domain.Task, error) {
	tasks, err := s.repo.GetAll(ctx, id, completed)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TodoService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, core_errors.ErrTaskNotFound
		}
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TodoService) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_errors.ErrTaskNotFound
	}

	return err
}

func (s *TodoService) CompleteTask(ctx context.Context, id int, completed *bool) (domain.Task, error) {
	task, err := s.repo.Complete(ctx, id, completed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, core_errors.ErrTaskNotFound
		}
		return domain.Task{}, err
	}

	return task, nil
}
