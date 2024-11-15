package drivers

import (
    "github.com/jmoiron/sqlx"
    "database/sql"
    "github.com/pressly/goose/v3"
    _ "github.com/mattn/go-sqlite3"
)

func Connect(connURL string) (*sqlx.DB, error) {
    db, err := sql.Open("sqlite3", connURL) 
    if err != nil {
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    if err := Migrate(db); err != nil {
        return nil, err
    }

    return sqlx.NewDb(db, "sqlite3"), nil
}

func Migrate(db *sql.DB) error {
    return  goose.Up(db, "internal/db/migrations") 
}

