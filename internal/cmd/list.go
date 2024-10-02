package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ncfex/tasks/internal/task"
	"github.com/ncfex/tasks/internal/utils"
	"github.com/spf13/cobra"
)

var (
	showAll         bool
	selectedColumns []string
	defaultColumns  = []string{"id", "description", "created"}
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
	filter := &task.TaskFilter{
		IncludeCompleted: showAll,
	}

	selector := task.NewTaskSelector(
		task.TaskFieldID,
		task.TaskFieldDescription,
		task.TaskFieldCreatedAt,
		task.TaskFieldIsCompleted,
	)

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

	headers := defaultColumns
	fmt.Fprintln(writer, joinWithTabs(headers))

	for _, t := range tasks {
		row := []string{
			fmt.Sprintf("%d", t.ID),
			t.Description,
			utils.HumanReadableTime(t.CreatedAt),
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
