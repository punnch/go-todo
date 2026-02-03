package todo

import (
	"time"
)

// todo:
/*
1. Realize why do I need a TaskUseCase
*/

// Entity
type Task struct {
	ID         int
	Title      string
	Descripton string
	Completed  bool
	CreatedAt  time.Time
}

func NewTask(title string, description string) Task {
	return Task{
		Title:      title,
		Descripton: description,
		Completed:  false,
		CreatedAt:  time.Now(),
	}
}

func (t *Task) Complete() {
	t.Completed = true
}

type TaskUseCase interface {
	Complete()
}
