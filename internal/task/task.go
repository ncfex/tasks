package task

import "time"

type TaskField string

const (
	TaskFieldID          TaskField = "id"
	TaskFieldDescription TaskField = "description"
	TaskFieldIsCompleted TaskField = "is_completed"
	TaskFieldCreatedAt   TaskField = "created_at"
)

type Task struct {
	ID          int
	Description string
	IsCompleted bool
	CreatedAt   time.Time
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

func NewTask(description string) Task {
	return Task{
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
}
