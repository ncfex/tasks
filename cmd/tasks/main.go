package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ncfex/tasks/internal/cli"
	"github.com/ncfex/tasks/internal/storage/csv"
	"github.com/ncfex/tasks/internal/task"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	storageDir := filepath.Join(homeDir, ".tasks")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	storagePath := filepath.Join(storageDir, "tasks.csv")
	repository := csv.NewRepository(storagePath)

	service := task.NewService(repository)

	app := cli.NewApp(service)
	if err := app.Run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
