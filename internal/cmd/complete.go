package cmd

import (
	"fmt"
	"strconv"

	"github.com/ncfex/tasks/internal/task"
	"github.com/spf13/cobra"
)

var CompleteCommand = &cobra.Command{
	Use:   "complete",
	Short: "Complete a TODO",
	Long:  "Complete a TODO from your list",
	Run: func(cC *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a ID of a TODO")
			return
		}

		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error converting string to int")
			return
		}

		err = task.CompleteTask(taskID)
		if err != nil {
			fmt.Printf("Error completing TODO: %v\n", err)
			return
		}

		fmt.Println("TODO completed succesfully.")
	},
}
