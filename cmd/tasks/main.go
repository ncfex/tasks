package main

import (
	"log"
	"os"

	"github.com/ncfex/tasks/internal/cli"
)

func main() {
	app := cli.NewApp()
	if err := app.Run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
