package main

import (
	"fmt"
	"time"

	"github.com/ncfex/todo-app/internal/task"
	"github.com/ncfex/todo-app/internal/utils"
)

func main() {

	newTask := task.Task{
		Opened: time.Now(),
	}

	fmt.Println(utils.HumanReadableTime(newTask.Opened))
}
