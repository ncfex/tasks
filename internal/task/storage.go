package task

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

const filename = "todos.csv"

func SaveTask(t Task) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	tasks, err := GetAllTasks()
	if err != nil {
		return err
	}

	t.ID = len(tasks) + 1

	record := []string{
		strconv.Itoa(t.ID),
		t.Description,
		strconv.FormatBool(t.Completed),
		t.CreatedAt.Format(time.RFC3339),
	}

	return writer.Write(record)
}

func GetAllTasks() ([]Task, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

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
