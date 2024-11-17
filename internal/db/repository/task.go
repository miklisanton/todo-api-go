package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"todo-api/internal/db/models"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type ITaskRepo interface {
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) ([]models.Task, error)
	Delete(ctx context.Context, id int) error
	GetTasksAfterDue(ctx context.Context) ([]models.Task, error)
}

type TaskRepo struct {
	db *sqlx.DB
}

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrNoTitle       = errors.New("title is required")
	ErrAlreadyExists = errors.New("task with given id already exists")
)

func NewTaskRepo(db *sqlx.DB) ITaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) Create(ctx context.Context, task *models.Task) error {
	query := `
    INSERT INTO task(id, title, description, due_date)
    VALUES($1, $2, $3, $4)
    RETURNING id, title, description, due_date, completed, overdue
    `
	row := r.db.QueryRowxContext(ctx, query, task.ID, task.Title, task.Description, task.DueDate)
	err := row.StructScan(task)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrAlreadyExists
			} else if sqliteErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
				return ErrNoTitle
			}
		}
		return err
	}
	return nil
}

func (r *TaskRepo) Update(ctx context.Context, task *models.Task) error {
	// Build query
	query := `UPDATE task SET`
	if task.Title != nil {
		query += " title = :title,"
	}
	if task.Description != nil {
		query += " description = :description,"
	}
	if task.DueDate != nil {
		query += " due_date = :due_date,"
	}
	if task.Completed != nil {
		query += " completed = :completed,"
	}
	if task.Overdue != nil {
		query += " overdue = :overdue,"
	}
	query = strings.TrimSuffix(query, ",") +
		" WHERE id = :id" +
		" RETURNING id, title, description, due_date, completed, overdue"

	rows, err := r.db.NamedQueryContext(ctx, query, task)
	if err != nil {
		return err
	}

	if rows.Next() {
		defer rows.Close()
		err = rows.StructScan(task)
	} else {
		return ErrTaskNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	task := &models.Task{}
	query := `SELECT * FROM task WHERE id = $1`
	err := r.db.GetContext(ctx, task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (r *TaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks := []models.Task{}
	query := `SELECT * FROM task`
	err := r.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM task WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepo) GetTasksAfterDue(ctx context.Context) ([]models.Task, error) {
	query := `SELECT * FROM task WHERE due_date < CURRENT_TIMESTAMP AND overdue = false`
	tasks := []models.Task{}
	err := r.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
