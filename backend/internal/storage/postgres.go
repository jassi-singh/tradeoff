package storage

import (
	"tradeoff/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresStore struct {
	DB *gorm.DB
}

func NewPostgresStore(config config.Config) (*PostgresStore, error) {
	// Configure GORM with PostgreSQL driver
	db, err := gorm.Open(postgres.Open(config.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Get the underlying SQL DB for ping check
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	store := &PostgresStore{DB: db}

	// Auto-migrate models
	if err := store.AutoMigrate(); err != nil {
		return nil, err
	}

	return store, nil
}

// AutoMigrate runs auto-migration for all models
func (s *PostgresStore) AutoMigrate() error {
	return s.DB.AutoMigrate(
		&PlayerModel{},
	)
}
