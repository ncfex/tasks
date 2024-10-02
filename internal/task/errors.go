package task

import (
	"errors"
	"fmt"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTask  = errors.New("invalid task")
)

type Error struct {
	Op  string
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}
