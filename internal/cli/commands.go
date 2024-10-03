package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

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
			return fmt.Sprintf("%d", t.ID)
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

func newAddCommand(service task.TaskService) *cobra.Command {
	var dueDateString string

	cmd := &cobra.Command{
		Use:   "add [description]",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(service, args[0], dueDateString)
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
		return fmt.Errorf("failed to create task: %w", err)
	}

	task, err := service.Create(description, dueDateTime)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Printf("Task created with ID: %d\n", task.ID)
	return nil
}

func newListCommand(service task.TaskService) *cobra.Command {
	var showAll bool
	var selectedColumns []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(service, showAll, selectedColumns)
		},
	}

	defaultColumns := []string{
		string(task.TaskFieldID),
		string(task.TaskFieldDescription),
		string(task.TaskFieldIsCompleted),
		string(task.TaskFieldCreatedAt),
		string(task.TaskFieldDueDate),
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks (including completed)")
	cmd.Flags().StringSliceVarP(&selectedColumns, "columns", "c", defaultColumns, "Columns to display")

	return cmd
}

func runList(service task.TaskService, showAll bool, selectedColumns []string) error {
	var displayColumns []Column
	var selectedFields []task.TaskField

	for _, colName := range selectedColumns {
		col, exists := columns[colName]
		if !exists {
			return fmt.Errorf("invalid column: %s", colName)
		}
		displayColumns = append(displayColumns, col)
		selectedFields = append(selectedFields, col.Field)
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

	var headers []string
	for _, col := range displayColumns {
		headers = append(headers, col.Header)
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, t := range tasks {
		var row []string
		for _, col := range displayColumns {
			row = append(row, col.Formatter(t))
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return nil
}

func newCompleteCommand(service task.TaskService) *cobra.Command {
	return &cobra.Command{
		Use:   "complete [task_id]",
		Short: "Mark a task as completed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %w", err)
			}

			if err := service.Complete(id); err != nil {
				return fmt.Errorf("failed to complete task: %w", err)
			}

			fmt.Printf("Task %d marked as completed\n", id)
			return nil
		},
	}
}
