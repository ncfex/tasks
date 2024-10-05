package config

import "github.com/ncfex/tasks/internal/task"

type ServiceMode string

const (
	ServiceModeSQL  ServiceMode = "sql"
	ServiceModeCSV  ServiceMode = "csv"
	ServiceModeJSON ServiceMode = "json"
)

type Config struct {
	filepath       string
	ServiceMode    ServiceMode      `json:"service_mode"`
	DisplayColumns []task.TaskField `json:"display_columns"`
}
