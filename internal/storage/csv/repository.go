package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

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
		return fmt.Errorf("failed to read tasks: %w", err)
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
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	var tasks []task.Task
	for _, record := range records {
		if len(record) != 5 {
			continue
		}

		id, _ := strconv.Atoi(record[0])
		isCompleted, _ := strconv.ParseBool(record[2])
		createdAt, _ := time.Parse(time.RFC3339, record[3])
		dueDate, _ := time.Parse(time.RFC3339, record[4])

		task := task.Task{
			ID:          id,
			Description: record[1],
			IsCompleted: isCompleted,
			CreatedAt:   createdAt,
			DueDate:     dueDate,
		}
		tasks = append(tasks, task)
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

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range tasks {
		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			strconv.FormatBool(t.IsCompleted),
			t.CreatedAt.Format(time.RFC3339),
			t.DueDate.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
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
		// TODO WRITE HEADERS
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
