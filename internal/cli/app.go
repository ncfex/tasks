package cli

import (
	"github.com/ncfex/tasks/internal/task"
	"github.com/spf13/cobra"
)

type App struct {
	rootCmd *cobra.Command
	service task.TaskService
}

func NewApp(service task.TaskService) *App {
	app := &App{
		service: service,
	}

	app.rootCmd = &cobra.Command{
		Use:   "tasks",
		Short: "Simple CLI todo app",
		Long:  "Simple CLI application for managing your todos",
	}

	app.setupCommands()
	return app
}

func (a *App) Run() error {
	return a.rootCmd.Execute()
}

func (a *App) setupCommands() {
	a.rootCmd.AddCommand(
		newAddCommand(a.service),
		newListCommand(a.service),
		newCompleteCommand(a.service),
	)
}
