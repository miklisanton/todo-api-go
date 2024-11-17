package services

import (
	"context"
	"time"
	"todo-api/internal/db/models"
	"todo-api/internal/db/repository"

	"github.com/rs/zerolog/log"
)

type ITaskService interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, id int) (*models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	UpdateOverdue(ctx context.Context) error
	UpdateTask(ctx context.Context, task *models.Task) error
	SetCompleted(ctx context.Context, id int, completed bool) (*models.Task, error)
	SetOverdue(ctx context.Context, id int, overdue bool) error
	DeleteTask(ctx context.Context, id int) error
}

type TaskService struct {
	Repo repository.ITaskRepo
}

func NewTaskService(taskRepo repository.ITaskRepo) ITaskService {
	return TaskService{taskRepo}
}

func (s TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	err := s.Repo.Create(ctx, task)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to create task")
		return err
	}
	return nil
}

func (s TaskService) GetTask(ctx context.Context, id int) (*models.Task, error) {
	task, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to get task with id %d", id)
		return nil, err
	}
	return task, nil
}

func (s TaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.Repo.GetAll(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to get tasks")
		return nil, err
	}
	return tasks, nil
}

func (s TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	// Check if task is overdue
	if task.DueDate != nil {
		var overdue bool
		if task.DueDate.Before(time.Now()) {
			overdue = true
		} else {
			overdue = false
		}
		task.Overdue = &overdue
	}

	err := s.Repo.Update(ctx, task)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to update task with id %d", *task.ID)
		return err
	}
	return nil
}

func (s TaskService) SetCompleted(ctx context.Context, id int, completed bool) (*models.Task, error) {
	task := &models.Task{
		ID:        &id,
		Completed: &completed,
	}
	err := s.Repo.Update(ctx, task)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to set completed task with id %d", id)
		return nil, err
	}
	return task, nil
}

func (s TaskService) SetOverdue(ctx context.Context, id int, overdue bool) error {
	task := &models.Task{
		ID:      &id,
		Overdue: &overdue,
	}
	err := s.Repo.Update(ctx, task)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to set overdue task with id %d", id)
		return err
	}
	return nil
}

func (s TaskService) DeleteTask(ctx context.Context, id int) error {
	err := s.Repo.Delete(ctx, id)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to delete task with id %d", id)
		return err
	}
	return nil
}

// UpdateOverdue fetches all overdue tasks and sets the overdue flag to true
func (s TaskService) UpdateOverdue(ctx context.Context) error {
	tasks, err := s.Repo.GetTasksAfterDue(ctx)
	log.Logger.Info().Msgf("found overdue tasks: %d", len(tasks))
	if err != nil {
		log.Logger.Error().Err(err).Msgf("failed to get tasks by due date")
		return err
	}
	for _, task := range tasks {
		err := s.SetOverdue(ctx, *task.ID, true)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("failed to set overdue for task with id %d", *task.ID)
			return err
		}
		log.Logger.Info().Msgf("task with id %d is overdue", *task.ID)
	}
	log.Logger.Info().Msg("overdue tasks updated")
	return nil
}
