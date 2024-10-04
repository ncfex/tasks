package task

import "github.com/google/uuid"

type Repository interface {
	Save(*Task) error
	GetByID(id uuid.UUID) (*Task, error)
	GetTaskByPartialId(id string) (*Task, error)
	List(*TaskSelector, *TaskFilter) ([]Task, error)
	Update(*Task) error
	Delete(*Task) error
}
