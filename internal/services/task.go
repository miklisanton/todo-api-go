package services

import (
	"context"
	"todo-api/internal/db/models"
	"todo-api/internal/db/repository"

	"github.com/rs/zerolog/log"
)

type ITaskService interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, id int) (*models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
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
	// Update task
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
