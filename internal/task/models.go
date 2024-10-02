package task

import (
	"errors"
	"time"
)

type TaskField string

const (
	TaskFieldID          TaskField = "id"
	TaskFieldDescription TaskField = "description"
	TaskFieldIsCompleted TaskField = "is_completed"
	TaskFieldCreatedAt   TaskField = "created_at"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskSelector struct {
	Fields map[TaskField]bool
}

type TaskFilter struct {
	IncludeCompleted bool
}

func NewTaskSelector(fields ...TaskField) *TaskSelector {
	selector := &TaskSelector{
		Fields: make(map[TaskField]bool),
	}

	if len(fields) == 0 {
		selector.Fields[TaskFieldID] = true
		selector.Fields[TaskFieldDescription] = true
		selector.Fields[TaskFieldIsCompleted] = true
		selector.Fields[TaskFieldCreatedAt] = true
		return selector
	}

	for _, field := range fields {
		selector.Fields[field] = true
	}
	return selector
}

func NewTaskFilter() *TaskFilter {
	return &TaskFilter{
		IncludeCompleted: false,
	}
}

func (t *Task) Validate() error {
	if t.Description == "" {
		return errors.New("task description cannot be empty")
	}
	return nil
}
