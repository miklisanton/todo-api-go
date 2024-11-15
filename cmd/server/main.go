package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/db/drivers"
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
	db, err = drivers.Connect(cfg.Db.Name)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected")
}

func main() {
	// Setup services
	// Setup controllers
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
	pg := e.Group("/api1/public")

	// Endpoints
	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
