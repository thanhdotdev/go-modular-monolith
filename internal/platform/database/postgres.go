package database

import (
	"fmt"
	"project-example/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
	if !cfg.Enabled() {
		return nil, fmt.Errorf("database dsn is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
