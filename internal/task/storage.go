package task

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

const filepath = "todos.csv"

func loadFile(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to lock file: %w", err)
	}

	return file, nil
}

func closeFile(file *os.File) error {
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		return fmt.Errorf("failed to unlock file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func SaveTask(t Task) error {
	tasks, err := GetAllTasks(NewTaskSelector(), &TaskFilter{})
	if err != nil {
		return err
	}
	t.ID = len(tasks) + 1

	file, err := loadFile(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeFile(file); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = file.Seek(0, 2)
	if err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(t.ID),
		t.Description,
		strconv.FormatBool(t.IsCompleted),
		t.CreatedAt.Format(time.RFC3339),
	}

	return writer.Write(record)
}

func GetTaskById(id int) (Task, error) {
	tasks, err := GetAllTasks(NewTaskSelector(), &TaskFilter{})
	if err != nil {
		return Task{}, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, fmt.Errorf("No Task found with this ID: %d", id)
}

func CompleteTask(id int) error {
	t, err := GetTaskById(id)
	if err != nil {
		return err
	}

	tasks, err := GetAllTasks(NewTaskSelector(), &TaskFilter{})
	if err != nil {
		return err
	}

	file, err := loadFile(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeFile(file); cerr != nil && err == nil {
			err = cerr
		}
	}()

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == t.ID {
			tasks[i].IsCompleted = true
		}
	}

	var records [][]string
	for _, t := range tasks {
		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			strconv.FormatBool(t.IsCompleted),
			t.CreatedAt.Format(time.RFC3339),
		}

		records = append(records, record)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func GetAllTasks(selector *TaskSelector, filter *TaskFilter) ([]Task, error) {
	if selector == nil {
		selector = NewTaskSelector()
	}
	if filter == nil {
		filter = NewTaskFilter()
	}

	file, err := loadFile(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := closeFile(file); cerr != nil && err == nil {
			err = cerr
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, record := range records {
		task := Task{}

		if selector.Fields[TaskFieldID] {
			task.ID, _ = strconv.Atoi(record[0])
		}

		if selector.Fields[TaskFieldDescription] {
			task.Description = record[1]
		}

		isCompleted, _ := strconv.ParseBool(record[2])
		if !filter.IncludeCompleted && isCompleted {
			continue
		}
		if selector.Fields[TaskFieldIsCompleted] {
			task.IsCompleted = isCompleted
		}

		if selector.Fields[TaskFieldCreatedAt] {
			task.CreatedAt, _ = time.Parse(time.RFC3339, record[3])
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
