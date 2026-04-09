//go:build integration

package orderpostgres

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	orderdomain "project-example/internal/modules/order/domain"
	"project-example/internal/platform/config"
	"project-example/internal/platform/database"
	"runtime"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
)

func TestRepositoryFindByID(t *testing.T) {
	db, cfg, cleanup := openTestPostgres(t)
	defer cleanup()

	if err := database.MigrateUp(cfg); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	repo := NewRepository(db)

	order, err := repo.FindByID(context.Background(), "ord-001")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if order.ID != "ord-001" {
		t.Fatalf("expected order id ord-001, got %s", order.ID)
	}
}

func TestRepositoryFindByIDReturnsDomainErrorWhenMissing(t *testing.T) {
	db, cfg, cleanup := openTestPostgres(t)
	defer cleanup()

	if err := database.MigrateUp(cfg); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	repo := NewRepository(db)

	_, err := repo.FindByID(context.Background(), "ord-404")
	if !errors.Is(err, orderdomain.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func openTestPostgres(t *testing.T) (*gorm.DB, config.DatabaseConfig, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		t.Skipf("docker daemon not available: %v", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_DB=project_example_test",
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
		},
	}, func(hostConfig *docker.HostConfig) {
		hostConfig.AutoRemove = true
		hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	resource.Expire(120)

	dsn := fmt.Sprintf(
		"host=localhost user=postgres password=postgres dbname=project_example_test port=%s sslmode=disable",
		resource.GetPort("5432/tcp"),
	)
	cfg := config.DatabaseConfig{
		DSN:           dsn,
		MigrationsDir: migrationsDir(t),
	}

	var db *gorm.DB
	if err := pool.Retry(func() error {
		var openErr error
		db, openErr = database.OpenPostgres(cfg)
		if openErr != nil {
			return openErr
		}

		sqlDB, openErr := db.DB()
		if openErr != nil {
			return openErr
		}

		return sqlDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		t.Fatalf("connect postgres container: %v", err)
	}

	cleanup := func() {
		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}

		_ = pool.Purge(resource)
	}

	return db, cfg, cleanup
}

func migrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve current test file path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", "..", "migrations"))
}
