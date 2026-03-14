package todo

import (
	"context"

	"github.com/punnch/go-todo/internal/core/domains/todo"
)

// contract to service
type TodoRepository interface {
	Create(ctx context.Context, task todo.Task) (todo.Task, error)
	GetAll(ctx context.Context, id *int, completed *bool) ([]todo.Task, error)
	Get(ctx context.Context, id int) (todo.Task, error)
	Delete(ctx context.Context, id int) error
	Complete(ctx context.Context, id int, completed bool) (todo.Task, error)
}
