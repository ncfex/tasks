package task

import "time"

type Task struct {
	ID          int
	Description string
	Completed   bool
	CreatedAt   time.Time
}

func NewTask(description string) Task {
	return Task{
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}
