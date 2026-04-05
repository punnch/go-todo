package domain

import (
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTask(
	title string,
	description string,
	completed bool,
	createdAt time.Time,
) Task {
	return Task{
		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   createdAt,
	}
}

func CreateTask(
	title string,
	description string,
) Task {
	var (
		completed = false
		createdAt = time.Now()
	)

	return NewTask(
		title,
		description,
		completed,
		createdAt,
	)
}
