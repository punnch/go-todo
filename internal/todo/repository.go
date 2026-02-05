package todo

import (
	"context"
	"errors"
)

// contract to service

type TodoRepository interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetAll(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id int) (Task, error)
	Delete(ctx context.Context, id int) error
	Complete(ctx context.Context, id int) (Task, error)
}

var (
	ErrNotFound     error = errors.New("task not found")
	ErrInvalidTitle error = errors.New("invalid task title")
)
