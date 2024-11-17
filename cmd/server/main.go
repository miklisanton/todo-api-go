package main

import (
	"os"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/db/drivers"
	"todo-api/internal/db/repository"
	"todo-api/internal/handlers"
	"todo-api/internal/jobs"
	"todo-api/internal/requests"
	"todo-api/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	db  *sqlx.DB
	cfg *config.Config
)

func init() {
	// Setup logger
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Info().Msg("Logger initialized")
	cfgPath, err := config.ParseCLI()
	// Parse CLI arguments
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse CLI")
	}
	// Read config
	cfg, err = config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}
	log.Info().Msg("Config loaded")
	// Connect to database
	db, err = drivers.Connect(cfg.Db.Name, "internal/db/migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected")
}

func main() {
	// Setup repositories
	taskRepo := repository.NewTaskRepo(db)
	// Setup services
	taskService := services.NewTaskService(taskRepo)
	// Setup controllers
	taskController := handlers.NewTaskController(taskService, time.Duration(cfg.Server.Timeout)*time.Second)
	// Setup echo
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))
	e.Validator = &requests.CustomValidator{Validator: validator.New()}
	pg := e.Group("/api1/public")

	// Endpoints
	pg.POST("/tasks", taskController.CreateTask)
	pg.GET("/tasks/:id", taskController.GetTask)
	pg.GET("/tasks", taskController.GetTasks)
	pg.PATCH("/tasks/:id/completed", taskController.SetCompleted)
	pg.PUT("/tasks/:id", taskController.UpdateTask)
	// Start overdue tasks monitor
	exitChan := make(chan os.Signal, 1)
	dateWorker := jobs.NewDateWorker(taskService, exitChan)
	go dateWorker.MonitorDueDate(time.Duration(10) * time.Second)
	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
