package todo

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// contract to service

type TodoRepository interface {
	Create(ctx context.Context, conn *pgx.Conn, task Task) error
}
