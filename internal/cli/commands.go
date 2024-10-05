package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ncfex/tasks/internal/config"
	"github.com/ncfex/tasks/internal/task"
	"github.com/ncfex/tasks/internal/utils"
	"github.com/spf13/cobra"
)

type Column struct {
	Header    string
	Field     task.TaskField
	Formatter func(t task.Task) string
}

var columns = map[string]Column{
	string(task.TaskFieldID): {
		Header: strings.ToUpper(string(task.TaskFieldID)),
		Field:  task.TaskFieldID,
		Formatter: func(t task.Task) string {
			return t.ID.String()[0:8]
		},
	},
	string(task.TaskFieldDescription): {
		Header: strings.ToUpper(string(task.TaskFieldDescription)),
		Field:  task.TaskFieldDescription,
		Formatter: func(t task.Task) string {
			return t.Description
		},
	},
	string(task.TaskFieldIsCompleted): {
		Header: strings.ToUpper(string(task.TaskFieldIsCompleted)),
		Field:  task.TaskFieldIsCompleted,
		Formatter: func(t task.Task) string {
			if t.IsCompleted {
				return "OK"
			}
			return "-"
		},
	},
	string(task.TaskFieldCreatedAt): {
		Header: strings.ToUpper(string(task.TaskFieldCreatedAt)),
		Field:  task.TaskFieldCreatedAt,
		Formatter: func(t task.Task) string {
			return utils.FormatTimeToHuman(t.CreatedAt)
		},
	},
	string(task.TaskFieldDueDate): {
		Header: strings.ToUpper(string(task.TaskFieldDueDate)),
		Field:  task.TaskFieldDueDate,
		Formatter: func(t task.Task) string {
			return utils.FormatTimeToHuman(t.DueDate)
		},
	},
}

var defaultColumns = []task.TaskField{
	task.TaskFieldID,
	task.TaskFieldDescription,
	task.TaskFieldCreatedAt,
	task.TaskFieldDueDate,
	task.TaskFieldIsCompleted,
}

func newAddCommand(a *App) *cobra.Command {
	var dueDateString string

	cmd := &cobra.Command{
		Use:   "add [description]",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(a.service, args[0], dueDateString)
		},
	}

	cmd.Flags().StringVarP(&dueDateString, "due", "d", "", "Due date for the task")

	return cmd
}

func runAdd(service task.TaskService, description string, dueDate string) error {
	if dueDate == "" {
		dueDate = "tomorrow"
	}

	dueDateTime, err := utils.ParseHumanToTime(dueDate)
	if err != nil {
		return fmt.Errorf("failed to create parse date: %w", err)
	}

	task, err := service.Create(description, dueDateTime)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Printf("Task created with ID: %s\n", task.ID.String()[0:8])
	return nil
}

func newListCommand(a *App) *cobra.Command {
	var showAll bool
	var selectedColumns []string
	var saveColumns bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(a.cfg.DisplayColumns) == 0 {
				if err := a.cfg.UpdateDisplayColumns(defaultColumns); err != nil {
					return fmt.Errorf("failed to set default columns: %w", err)
				}
			}

			var columnsToUse []string
			if len(selectedColumns) > 0 {
				selectedTaskFields := make([]task.TaskField, len(selectedColumns))
				for i, col := range selectedColumns {
					selectedTaskFields[i] = task.TaskField(col)
				}

				if saveColumns {
					if err := a.cfg.UpdateDisplayColumns(selectedTaskFields); err != nil {
						return fmt.Errorf("failed to update display columns in config: %w", err)
					}
				}

				columnsToUse = selectedColumns
			} else {
				columnsToUse = make([]string, 0, len(a.cfg.DisplayColumns))
				for _, col := range a.cfg.DisplayColumns {
					columnsToUse = append(columnsToUse, string(col))
				}
			}

			return runList(a.service, showAll, columnsToUse)
		},
	}

	displayColumnsString := make([]string, 0, len(a.cfg.DisplayColumns))
	for _, c := range a.cfg.DisplayColumns {
		displayColumnsString = append(displayColumnsString, string(c))
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks (including completed)")
	cmd.Flags().StringSliceVarP(&selectedColumns, "columns", "c", displayColumnsString, "Columns to display")
	cmd.Flags().BoolVarP(&saveColumns, "save", "s", false, "Save selected columns to config")

	return cmd
}

func runList(service task.TaskService, showAll bool, selectedColumns []string) error {
	displayColumns := make([]Column, 0, len(selectedColumns))
	selectedFields := make([]task.TaskField, 0, len(selectedColumns))

	if len(selectedColumns) == 0 {
		displayColumns = make([]Column, 0, len(columns))
		selectedFields = make([]task.TaskField, 0, len(columns))

		for _, col := range columns {
			displayColumns = append(displayColumns, col)
			selectedFields = append(selectedFields, col.Field)
		}
	} else {
		for _, colName := range selectedColumns {
			col, exists := columns[colName]
			if !exists {
				return fmt.Errorf("invalid column: %s", colName)
			}
			displayColumns = append(displayColumns, col)
			selectedFields = append(selectedFields, col.Field)
		}
	}

	selector := task.NewTaskSelector(selectedFields...)
	filter := &task.TaskFilter{IncludeCompleted: showAll}

	tasks, err := service.List(selector, filter)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	headers := make([]string, 0, len(displayColumns))
	for _, col := range displayColumns {
		headers = append(headers, col.Header)
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	row := make([]string, 0, len(displayColumns))
	for _, t := range tasks {
		row = row[:0]
		for _, col := range displayColumns {
			row = append(row, col.Formatter(t))
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return nil
}

func newCompleteCommand(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "complete [task_id]",
		Short: "Mark a task as completed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			idString := args[0]
			if err := a.service.Complete(idString); err != nil {
				return fmt.Errorf("failed to complete task: %w", err)
			}

			fmt.Printf("Task %s marked as completed\n", idString)
			return nil
		},
	}
}

func newDeleteCommand(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [task_id]",
		Short: "Delete a task from TODOs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			idString := args[0]
			if err := a.service.Delete(idString); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			fmt.Printf("Task %s is deleted\n", idString)
			return nil
		},
	}
}

func newUpdateServiceModeCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-mode [mode]",
		Short: "Update default storage mode",
		Long:  "Update the default storage mode (sql, json, or csv) for the tasks application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := config.ServiceMode(args[0])
			switch mode {
			case config.ServiceModeSQL, config.ServiceModeJSON, config.ServiceModeCSV:
			default:
				return fmt.Errorf("invalid mode: %s. Must be one of: sql, json, csv", mode)
			}

			if err := a.cfg.UpdateServiceMode(mode); err != nil {
				return fmt.Errorf("failed to update service mode: %w", err)
			}

			fmt.Printf("Successfully updated default storage mode to: %s\n", mode)
			return nil
		},
	}

	return cmd
}
