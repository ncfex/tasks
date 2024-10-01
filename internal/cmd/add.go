package cmd

import (
	"fmt"

	"github.com/ncfex/todo-app/internal/task"
	"github.com/spf13/cobra"
)

var AddCommand = &cobra.Command{
	Use:   "Add",
	Short: "Add a new TODO",
	Long:  "Add a new todo item to your list",
	Run: func(cC *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a TODO description")
			return
		}

		description := args[0]
		todo := task.NewTask(description)
		err := task.SaveTask(todo)
		if err != nil {
			fmt.Printf("Error saving todo: %v\n", err)
			return
		}
		fmt.Println("Todo added successfully!")
	},
}
