package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// Migrator handles database migrations
type Migrator struct {
	migrationsPath string
	db             *sqlx.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sqlx.DB, migrationsPath string) (*Migrator, error) {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return &Migrator{
		db:             db,
		migrationsPath: absPath,
	}, nil
}

// Up applies all available migrations
func (m *Migrator) Up() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Down rolls back all migrations
func (m *Migrator) Down() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

// CreateMigration creates new migration files
func (m *Migrator) CreateMigration(name string) error {
	upFile := filepath.Join(m.migrationsPath, fmt.Sprintf("%s_up.sql", name))
	downFile := filepath.Join(m.migrationsPath, fmt.Sprintf("%s_down.sql", name))

	up, err := os.Create(upFile)
	if err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}
	defer up.Close()

	down, err := os.Create(downFile)
	if err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}
	defer down.Close()

	fmt.Printf("Created migration files: %s and %s\n", upFile, downFile)

	return nil
}
