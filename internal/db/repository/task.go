package repository

import (
	"github.com/jmoiron/sqlx"
	"todo-api/internal/models"
)

type ITaskRepo interface {
	Create(task *models.Task) error
	Update(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetAll() ([]models.Task, error)
	Delete(id int) error
}
