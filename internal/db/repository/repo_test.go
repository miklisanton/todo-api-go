package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"todo-api/internal/db/drivers"
	"todo-api/internal/db/models"
	"todo-api/internal/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var r ITaskRepo

func TestMain(m *testing.M) {
	// Setup logger
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Info().Msg("Logger initialized")
	// Setup database
	db, err := drivers.Connect("test.db", "../migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	r = NewTaskRepo(db)
	// Run tests
	m.Run()

	drivers.Down(db, "../migrations")
}

func TestCreate(t *testing.T) {
	// Create a task
	title := "task1"
	description := "Test description"

	task := models.Task{
		Title:       &title,
		Description: &description,
	}
	err := r.Create(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	} else {
		utils.PrintTask(task)
	}
}

func TestCreateError(t *testing.T) {
	// Create a task without title
	description := "Test description"

	task := models.Task{
		Description: &description,
	}
	err := r.Create(context.TODO(), &task)
	if err == nil {
		t.Errorf("Expected error creating task")
	} else {
		t.Logf("Error creating task: %v", err)
	}
}

func TestCreateError2(t *testing.T) {
	// Create a task with existing ID

	description := "Test description"
	title := "task1"
	id := 1

	task := models.Task{
		ID:          &id,
		Title:       &title,
		Description: &description,
	}
	err := r.Create(context.TODO(), &task)
	if err == nil {
		t.Errorf("Expected error creating task")
	} else {
		t.Logf("Error creating task: %v", err)
	}
}

func TestCreate2(t *testing.T) {
	// Create a task
	id := 2
	title := "task2"
	description := "Test description"

	task := models.Task{
		ID:          &id,
		Title:       &title,
		Description: &description,
	}
	err := r.Create(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	} else {
		utils.PrintTask(task)
	}
}

func TestCreate3(t *testing.T) {
	// Create a task
	title := "task3"

	task := models.Task{
		Title: &title,
	}
	err := r.Create(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	} else {
		utils.PrintTask(task)
	}
}

func TestUpdate(t *testing.T) {
	// Update a task
	title := "task1 updated"
	description := "Test description updated"

	id := 1
	task := models.Task{
		ID:          &id,
		Title:       &title,
		Description: &description,
	}
	err := r.Update(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error updating task: %v", err)
	} else {
		utils.PrintTask(task)
	}
}

func TestUpdate2(t *testing.T) {
	// Update a task
	id := 2
	completed := true
	task := models.Task{
		ID:        &id,
		Completed: &completed,
	}
	err := r.Update(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error updating task: %v", err)
	} else {
		utils.PrintTask(task)
	}
}

func TestUpdateError(t *testing.T) {
	// Update not existing task
	id := 999
	description := "Test description updated"
	task := models.Task{
		ID:          &id,
		Description: &description,
	}
	err := r.Update(context.TODO(), &task)
	if err == nil {
		t.Errorf("Expected error updating task")
	} else {
		t.Logf("Error updating task: %v", err)
	}
}

func TestGetAll(t *testing.T) {
	// Get all tasks
	tasks, err := r.GetAll(context.TODO())
	if err != nil {
		t.Errorf("Error getting tasks: %v", err)
	}
	for _, task := range tasks {
		utils.PrintTask(task)
		fmt.Println()
	}
}

func TestGetByID(t *testing.T) {
	// Get a task by ID
	id := 1
	task, err := r.GetByID(context.TODO(), id)
	if err != nil {
		t.Errorf("Error getting task: %v", err)
	} else {
		utils.PrintTask(*task)
	}
}

func TestGetByIDError(t *testing.T) {
	// Get a task by ID
	id := 999
	_, err := r.GetByID(context.TODO(), id)
	if err == nil {
		t.Errorf("Expected error getting task")
	} else {
		t.Logf("Error getting task: %v", err)
	}
}

func TestDelete(t *testing.T) {
	// Delete a task
	id := 1
	err := r.Delete(context.TODO(), id)
	if err != nil {
		t.Errorf("Error deleting task: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	// Delete a task
	id := 999
	err := r.Delete(context.TODO(), id)
	if err == nil {
		t.Errorf("Expected error deleting task")
	} else {
		t.Logf("Error deleting task: %v", err)
	}
}
