package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "simple CLI todo app",
	Long:  "simple CLI application for managing your todos",
}

func init() {
	RootCmd.AddCommand(AddCommand)
}
