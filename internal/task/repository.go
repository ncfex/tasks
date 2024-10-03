package task

type Repository interface {
	Save(*Task) error
	GetByID(int) (*Task, error)
	List(*TaskSelector, *TaskFilter) ([]Task, error)
	Update(*Task) error
	Delete(*Task) error
}
