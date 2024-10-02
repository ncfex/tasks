package task

import (
	"time"
)

type TaskService interface {
	Create(description string) (*Task, error)
	GetByID(id int) (*Task, error)
	List(selector *TaskSelector, filter *TaskFilter) ([]Task, error)
	Complete(id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) TaskService {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(description string) (*Task, error) {
	task := &Task{
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}

	if err := task.Validate(); err != nil {
		return nil, &Error{Op: "Create", Err: err}
	}

	if err := s.repository.Save(task); err != nil {
		return nil, &Error{Op: "Create", Err: err}
	}

	return task, nil
}

func (s *service) GetByID(id int) (*Task, error) {
	task, err := s.repository.GetByID(id)
	if err != nil {
		return nil, &Error{Op: "GetByID", Err: err}
	}
	return task, nil
}

func (s *service) List(selector *TaskSelector, filter *TaskFilter) ([]Task, error) {
	if selector == nil {
		selector = NewTaskSelector()
	}
	if filter == nil {
		filter = NewTaskFilter()
	}

	tasks, err := s.repository.List(selector, filter)
	if err != nil {
		return nil, &Error{Op: "List", Err: err}
	}
	return tasks, nil
}

func (s *service) Complete(id int) error {
	task, err := s.repository.GetByID(id)
	if err != nil {
		return &Error{Op: "Complete", Err: err}
	}

	task.IsCompleted = true
	if err := s.repository.Update(task); err != nil {
		return &Error{Op: "Complete", Err: err}
	}

	return nil
}
