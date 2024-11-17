package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"todo-api/internal/db/models"
	"todo-api/internal/db/repository"
	"todo-api/internal/requests"
	"todo-api/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TaskController struct {
	TaskService services.ITaskService
	Timeout     time.Duration
}

func NewTaskController(taskService services.ITaskService, timeout time.Duration) *TaskController {
	return &TaskController{taskService, timeout}
}

func (tc *TaskController) CreateTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()

	taskReq := requests.PostTaskRequest{}
	if err := c.Bind(&taskReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to bind task")
		return c.JSON(http.StatusBadRequest, "failed to parse JSON")
	}
	if err := c.Validate(taskReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to validate task")
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	// Create task
	task := models.Task{
		Title: taskReq.Title,
	}
	// Parse due date
	if taskReq.DueDate == nil || *taskReq.DueDate == "" {
		task.DueDate = nil
	} else {
		parsed, err := time.Parse("2006-01-02", *taskReq.DueDate)
		if err != nil {
			log.Logger.Error().Err(err).Msg("failed to parse due date")
			return c.JSON(http.StatusBadRequest, "invalid due date")
		}
		task.DueDate = &parsed
	}
	// Set description if present
	if taskReq.Description == nil {
		task.Description = nil
	} else {
		task.Description = taskReq.Description
	}

	if err := tc.TaskService.CreateTask(ctx, &task); err != nil {
		log.Logger.Error().Err(err).Msg("failed to create task")
		if err == repository.ErrNoTitle {
			return c.JSON(http.StatusBadRequest, "task title is required")
		} else if err == repository.ErrAlreadyExists {
			return c.JSON(http.StatusConflict, "task already exists")
		}
		return c.JSON(http.StatusInternalServerError, "failed to create task")
	}
	return c.JSON(http.StatusCreated, task)
}

func (tc *TaskController) GetTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	// Retrieve task id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to parse task id")
		return c.JSON(http.StatusBadRequest, "invalid task id")
	}

	task, err := tc.TaskService.GetTask(ctx, id)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to get task with id %d", id)
		if err == repository.ErrTaskNotFound {
			return c.JSON(http.StatusNotFound, "task not found")
		}
		return c.JSON(http.StatusInternalServerError, "failed to get task")
	}
	return c.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetTasks(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	tasks, err := tc.TaskService.GetTasks(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to get tasks")
		return c.JSON(http.StatusInternalServerError, "failed to get tasks")
	}
	return c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	// Retrieve task id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to parse task id")
		return c.JSON(http.StatusBadRequest, "invalid task id")
	}

	taskReq := requests.PutTaskRequest{}
	if err := c.Bind(&taskReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to bind task")
		return c.JSON(http.StatusBadRequest, "failed to parse JSON")
	}
	if err := c.Validate(taskReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to validate task")
		return c.JSON(http.StatusBadRequest, "invalid request")
	}
	// Parse due date
	var dueDate *time.Time

	if taskReq.DueDate == nil {
		dueDate = nil
	} else {
		parsed, err := time.Parse("2006-01-02", *taskReq.DueDate)
		dueDate = &parsed
		if err != nil {
			log.Logger.Error().Err(err).Msg("failed to parse due date")
			return c.JSON(http.StatusBadRequest, "invalid due date")
		}
	}
	// Update task
	task := models.Task{
		ID:          &id,
		Title:       taskReq.Title,
		Description: taskReq.Description,
		DueDate:     dueDate,
	}
	err = tc.TaskService.UpdateTask(ctx, &task)
	if err == repository.ErrTaskNotFound {
		// Create new task with provided ID
		if err := tc.TaskService.CreateTask(ctx, &task); err != nil {
			log.Logger.Error().Err(err).Msg("failed to create task")
			return c.JSON(http.StatusInternalServerError, "failed to create task")
		}
		return c.JSON(http.StatusCreated, task)
	} else if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to update task with id %d", id)
		return c.JSON(http.StatusInternalServerError, "failed to update task")
	}

	return c.JSON(http.StatusOK, task)
}

func (tc *TaskController) SetCompleted(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	// Retrieve task id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to parse task id")
		return c.JSON(http.StatusBadRequest, "invalid task id")
	}

	completedReq := requests.PatchTaskRequest{}
	if err := c.Bind(&completedReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to bind task")
		return c.JSON(http.StatusBadRequest, "failed to parse JSON")
	}
	if err := c.Validate(completedReq); err != nil {
		log.Logger.Error().Err(err).Msg("failed to validate task")
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	taskUpdated, err := tc.TaskService.SetCompleted(ctx, id, completedReq.Completed)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to set completed task with id %d", id)
		if err == repository.ErrTaskNotFound {
			return c.JSON(http.StatusNotFound, "task not found")
		}
		return c.JSON(http.StatusInternalServerError, "failed to set completed")
	}
	return c.JSON(http.StatusOK, taskUpdated)
}

func (tc *TaskController) DeleteTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	// Retrieve task id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to parse task id")
		return c.JSON(http.StatusBadRequest, "invalid task id")
	}

	err = tc.TaskService.DeleteTask(ctx, id)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to delete task with id %d", id)
		if err == repository.ErrTaskNotFound {
			return c.JSON(http.StatusNotFound, "task not found")
		}
		return c.JSON(http.StatusInternalServerError, "failed to delete task")
	}
	return c.JSON(http.StatusOK, "task deleted")
}
