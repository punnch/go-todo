package todo

import (
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
}

func NewTask(title string, description string) Task {
	return Task{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}
