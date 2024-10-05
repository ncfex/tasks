package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ncfex/tasks/internal/config"
	"github.com/ncfex/tasks/internal/storage/csv"
	"github.com/ncfex/tasks/internal/storage/json"
	"github.com/ncfex/tasks/internal/storage/sql"
	"github.com/ncfex/tasks/internal/task"
	"github.com/spf13/cobra"
)

type App struct {
	rootCmd *cobra.Command
	service task.TaskService
	format  string
	cfg     *config.Config
}

func NewApp() *App {
	app := &App{}
	app.cfg = &config.Config{}

	err := app.cfg.Load()
	if err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	app.rootCmd = &cobra.Command{
		Use:   "tasks",
		Short: "Simple CLI todo app",
		Long:  "Simple CLI application for managing your todos",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.initializeService()
		},
	}

	app.rootCmd.PersistentFlags().StringVarP(&app.format, "format", "m", string(app.cfg.ServiceMode), "Storage format (json or csv)")

	app.setupCommands()
	return app
}

func (a *App) initializeService() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	storageDir := filepath.Join(homeDir, ".tasks")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	var repository task.Repository
	switch a.format {
	case "json":
		storagePath := filepath.Join(storageDir, "tasks.json")
		repository = json.NewRepository(storagePath)
	case "csv":
		storagePath := filepath.Join(storageDir, "tasks.csv")
		repository = csv.NewRepository(storagePath)
	case "sql":
		dbURL := os.Getenv("DB_URL")
		if dbURL == "" {
			log.Fatalf("Failed to get DB string")
		}
		repository, err = sql.NewRepository(dbURL)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", a.format)
	}

	a.service = task.NewService(repository)
	return nil
}

func (a *App) Run() error {
	return a.rootCmd.Execute()
}

func (a *App) setupCommands() {
	a.rootCmd.AddCommand(
		newAddCommand(a),
		newListCommand(a),
		newCompleteCommand(a),
		newDeleteCommand(a),
		newUpdateServiceModeCommand(a),
	)
}
