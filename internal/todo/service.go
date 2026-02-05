package todo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// Usecase
// business logic

type TodoService struct {
	repo TodoRepository
}

func (s *TodoService) CreateTask(ctx context.Context, title, desciption string) (Task, error) {
	task := NewTask(title, desciption)

	dbTask, err := s.repo.Create(ctx, task)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}

	return dbTask, nil
}

func (s *TodoService) GetAllTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TodoService) GetTask(ctx context.Context, id int) (Task, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}

	return task, nil
}

func (s *TodoService) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func (s *TodoService) CompleteTask(ctx context.Context, id int) (Task, error) {
	task, err := s.repo.Complete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, ErrNotFound
	}

	return task, nil
}
