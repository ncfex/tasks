package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ncfex/tasks/internal/task"
)

type repository struct {
	filepath string
	mu       sync.Mutex
}

func NewRepository(filepath string) *repository {
	return &repository{
		filepath: filepath,
	}
}

func (r *repository) Save(t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readTasks()
	if err != nil {
		return err
	}

	t.ID = r.nextID(tasks)
	tasks = append(tasks, *t)

	return r.writeTasks(tasks)
}

func (r *repository) GetByID(id int) (*task.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	for _, t := range tasks {
		if t.ID == id {
			return &t, nil
		}
	}

	return nil, task.ErrTaskNotFound
}

func (r *repository) List(selector *task.TaskSelector, filter *task.TaskFilter) ([]task.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	var filtered []task.Task
	for _, t := range tasks {
		if !filter.IncludeCompleted && t.IsCompleted {
			continue
		}
		filtered = append(filtered, t)
	}

	return filtered, nil
}

func (r *repository) Update(t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	found := false
	for i, task := range tasks {
		if task.ID == t.ID {
			tasks[i] = *t
			found = true
			break
		}
	}

	if !found {
		return task.ErrTaskNotFound
	}

	return r.writeTasks(tasks)
}

func (r *repository) Delete(t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readTasks()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	found := false
	indexToRemove := -1
	for i, task := range tasks {
		if task.ID == t.ID {
			indexToRemove = i
			found = true
			break
		}
	}

	if !found {
		return task.ErrTaskNotFound
	}

	tasks = append(tasks[:indexToRemove], tasks[indexToRemove+1:]...)

	return r.writeTasks(tasks)
}

func (r *repository) readTasks() ([]task.Task, error) {
	if err := r.ensureFile(); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(r.filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}

	return tasks, nil
}

func (r *repository) writeTasks(tasks []task.Task) error {
	if err := r.ensureFile(); err != nil {
		return err
	}

	file, err := os.OpenFile(r.filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (r *repository) ensureFile() error {
	dir := filepath.Dir(r.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(r.filepath); os.IsNotExist(err) {
		file, err := os.Create(r.filepath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		_, err = file.Write([]byte("[]"))
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

func (r *repository) nextID(tasks []task.Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}
