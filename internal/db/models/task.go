package models

import (
	"time"
)

type Task struct {
	ID          *int       `json:"id" db:"id"`
	Title       *string    `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	Completed   *bool      `json:"completed" db:"completed"`
	Overdue     *bool      `json:"overdue" db:"overdue"`
}
