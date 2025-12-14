package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/antonioparicio/numismaticapp/internal/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type MigrationService struct {
	db *sql.DB
}

func NewMigrationService(db *sql.DB) *MigrationService {
	return &MigrationService{
		db: db,
	}
}

func (s *MigrationService) RunMigrations() error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	d, err := iofs.New(migrations.FS, "sql")
	if err != nil {
		return fmt.Errorf("could not create iofs source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs", d,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
