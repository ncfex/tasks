package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ncfex/tasks/internal/storage/sql/database"
	"github.com/ncfex/tasks/internal/task"
)

type repository struct {
	db *database.Queries
}

func NewRepository(dbURL string) (task.Repository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	queries := database.New(db)
	return &repository{db: queries}, nil
}

func (r *repository) toSQLTask(t *task.Task) database.Task {
	return database.Task{
		ID:          uuid.New(),
		Description: t.Description,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		DueDate:     t.DueDate,
	}
}

func (r *repository) toDomainTask(t database.Task) task.Task {
	return task.Task{
		ID:          t.ID,
		Description: t.Description,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		DueDate:     t.DueDate,
	}
}

func (r *repository) Save(t *task.Task) error {
	sqlTask := r.toSQLTask(t)
	params := database.CreateTaskParams{
		Description: sqlTask.Description,
		DueDate:     sqlTask.DueDate,
	}

	_, err := r.db.CreateTask(context.Background(), params)
	return err
}

func (r *repository) GetByID(uuid uuid.UUID) (*task.Task, error) {
	sqlTask, err := r.db.GetTaskById(context.Background(), uuid)
	if err != nil {
		return nil, err
	}

	domainTask := r.toDomainTask(sqlTask)
	return &domainTask, nil
}

func (r *repository) GetTaskByPartialId(uuid string) (*task.Task, error) {
	nullUUID := sql.NullString{
		String: uuid,
		Valid:  true,
	}

	sqlTask, err := r.db.GetTaskByPartialId(context.Background(), nullUUID)
	if err != nil {
		return nil, err
	}

	domainTask := r.toDomainTask(sqlTask)
	return &domainTask, nil
}

func (r *repository) List(selector *task.TaskSelector, filter *task.TaskFilter) ([]task.Task, error) {
	sqlTasks, err := r.db.GetAllTasks(context.Background())
	if err != nil {
		return nil, err
	}

	tasks := make([]task.Task, len(sqlTasks))
	for i, sqlTask := range sqlTasks {
		tasks[i] = r.toDomainTask(sqlTask)
	}
	return tasks, nil
}

func (r *repository) Update(t *task.Task) error {
	sqlTask := r.toSQLTask(t)
	_, err := r.db.CompleteTask(context.Background(), sqlTask.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(t *task.Task) error {
	return r.db.DeleteTask(context.Background(), t.ID)
}
