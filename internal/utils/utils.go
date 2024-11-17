package utils

import (
	"fmt"
	"todo-api/internal/db/models"
)

func PrintTask(task models.Task) {
	fmt.Printf("ID: %d\n", *task.ID)
	fmt.Printf("Title: %s\n", *task.Title)
	if task.Description != nil {
		fmt.Printf("Description: %s\n", *task.Description)
	}
	if task.DueDate != nil {
		fmt.Printf("Due Date: %s\n", task.DueDate.Format("2006-01-02"))
	}
	if task.Completed != nil {
		fmt.Printf("Completed: %t\n", *task.Completed)
	}
	if task.Overdue != nil {
		fmt.Printf("Overdue: %t\n", *task.Overdue)
	}
}
