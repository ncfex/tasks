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
	tasks, err := GetAllTasks()
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
		strconv.FormatBool(t.Completed),
		t.CreatedAt.Format(time.RFC3339),
	}

	return writer.Write(record)
}

func GetAllTasks() ([]Task, error) {
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
		id, _ := strconv.Atoi(record[0])
		completed, _ := strconv.ParseBool(record[2])
		createdAt, _ := time.Parse(time.RFC3339, record[3])

		task := Task{
			ID:          id,
			Description: record[1],
			Completed:   completed,
			CreatedAt:   createdAt,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
