package config

type ServiceMode string

const (
	ServiceModeSQL  ServiceMode = "sql"
	ServiceModeCSV  ServiceMode = "csv"
	ServiceModeJSON ServiceMode = "json"
)

type Config struct {
	filepath    string
	ServiceMode ServiceMode `json:"service_mode"`
}
