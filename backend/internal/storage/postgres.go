package storage

import (
	"database/sql"
	"tradeoff/backend/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(config config.Config) (*PostgresStore, error) {
	db, err := sql.Open("pgx", config.Database.URL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{DB: db}, nil
}