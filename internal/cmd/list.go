package cmd

import (
	"fmt"

	"github.com/ncfex/tasks/internal/task"
	"github.com/spf13/cobra"
)

var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Long:  `Display a list of all your todos`,
	Run: func(cC *cobra.Command, args []string) {
		tasks, err := task.GetAllTasks()
		if err != nil {
			fmt.Printf("Error retrieving todos: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No todos found.")
			return
		}

		for _, task := range tasks {
			status := "[ ]"
			if task.Completed {
				status = "[x]"
			}
			fmt.Printf("%s %d: %s\n", status, task.ID, task.Description)
		}
	},
}
