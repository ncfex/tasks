package cmd

import (
	"fmt"
	"os"
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
	string(task.TaskFieldCreatedAt): {
		Header: strings.ToUpper(string(task.TaskFieldCreatedAt)),
		Field:  task.TaskFieldCreatedAt,
		Formatter: func(t task.Task) string {
			return utils.HumanReadableTime(t.CreatedAt)
		},
	},
	string(task.TaskFieldIsCompleted): {
		Header: strings.ToUpper(string(task.TaskFieldIsCompleted)),
		Field:  task.TaskFieldIsCompleted,
		Formatter: func(t task.Task) string {
			if t.IsCompleted {
				return "DONE"
			}
			return "-"
		},
	},
}

var (
	showAll         bool
	selectedColumns []string
	defaultColumns  = []string{
		string(task.TaskFieldID),
		string(task.TaskFieldDescription),
		string(task.TaskFieldCreatedAt),
		string(task.TaskFieldIsCompleted),
	}
)

var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Long:  `Display a list of all your todos`,
	RunE:  runList,
}

func init() {
	ListCommand.Flags().BoolVarP(&showAll, "show-all", "A", false, "Show all tasks (including completed ones)")
	ListCommand.Flags().StringSliceVarP(&selectedColumns, "columns", "c", defaultColumns, "Specify columns to display")
}

func runList(cmd *cobra.Command, args []string) error {
	var displayColumns []Column
	var selectedFields []task.TaskField

	selectedFields = append(selectedFields, task.TaskFieldIsCompleted)

	for _, colName := range selectedColumns {
		col, exists := columns[colName]
		if !exists {
			return fmt.Errorf("invalid column: %s", colName)
		}
		displayColumns = append(displayColumns, col)
		selectedFields = append(selectedFields, col.Field)
	}

	filter := &task.TaskFilter{
		IncludeCompleted: showAll,
	}

	selector := task.NewTaskSelector(selectedFields...)

	tasks, err := task.GetAllTasks(selector, filter)
	if err != nil {
		return fmt.Errorf("error retrieving todos: %v", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No todos found.")
		return nil
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
	defer writer.Flush()

	var headers []string
	for _, col := range displayColumns {
		headers = append(headers, col.Header)
	}
	fmt.Fprintln(writer, joinWithTabs(headers))

	for _, t := range tasks {
		var row []string
		for _, col := range displayColumns {
			row = append(row, col.Formatter(t))
		}
		fmt.Fprintln(writer, joinWithTabs(row))
	}

	return nil
}

func joinWithTabs(strings []string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += "\t"
		}
		result += s
	}
	return result
}
