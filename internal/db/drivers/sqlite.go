package drivers

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func Connect(connURL string, migrationsPath string) (*sqlx.DB, error) {
	log.Info().Msg("Connecting to database")
	db, err := sql.Open("sqlite3", connURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to database")
	if err := Up(db, migrationsPath); err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, "sqlite3"), nil
}

func Up(db *sql.DB, path string) error {
	goose.SetDialect("sqlite3")
	return goose.Up(db, path)
}

func Down(db *sqlx.DB, path string) error {
	goose.SetDialect("sqlite3")
	var sqlDB *sql.DB = db.DB
	return goose.Down(sqlDB, path)
}
