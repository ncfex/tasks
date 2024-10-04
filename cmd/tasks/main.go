package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ncfex/tasks/internal/cli"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	if err := app.Run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
