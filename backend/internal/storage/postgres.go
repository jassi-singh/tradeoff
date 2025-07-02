package storage

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{DB: db}, nil
}