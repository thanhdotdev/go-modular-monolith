package database

import (
	"database/sql"
	"errors"
	"path/filepath"
	"project-example/internal/platform/config"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MigrateUp(cfg config.DatabaseConfig) error {
	return runMigration(cfg, (*migrate.Migrate).Up)
}

func MigrateDown(cfg config.DatabaseConfig) error {
	return runMigration(cfg, (*migrate.Migrate).Down)
}

func runMigration(cfg config.DatabaseConfig, execute func(m *migrate.Migrate) error) error {
	m, sqlDB, err := newMigrator(cfg)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = m.Close()
		_ = sqlDB.Close()
	}()

	return ignoreNoChange(execute(m))
}

func newMigrator(cfg config.DatabaseConfig) (*migrate.Migrate, *sql.DB, error) {
	sqlDB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	absoluteDir, err := filepath.Abs(cfg.MigrationsDir)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+absoluteDir, "postgres", driver)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	return m, sqlDB, nil
}

func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
