package infrastructure

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MigrationService struct {
	db *pgxpool.Pool
}

func NewMigrationService(db *pgxpool.Pool) *MigrationService {
	return &MigrationService{db: db}
}

func (s *MigrationService) RunMigrations(ctx context.Context) error {
	// 1. Ensure schema_migrations table exists
	if err := s.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	// 2. Get applied migrations
	applied, err := s.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// 3. Get pending migrations from embedded FS
	pending, err := s.getPendingMigrations(applied)
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pending) == 0 {
		fmt.Println("No pending migrations.")
		return nil
	}

	// 4. Apply pending migrations
	for _, m := range pending {
		fmt.Printf("Applying migration: %s\n", m.Name)
		if err := s.applyMigration(ctx, m); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", m.Name, err)
		}
	}

	fmt.Println("All migrations applied successfully.")
	return nil
}

func (s *MigrationService) ensureMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			dirty BOOLEAN NOT NULL DEFAULT false,
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := s.db.Exec(ctx, query)
	return err
}

func (s *MigrationService) getAppliedMigrations(ctx context.Context) (map[int64]bool, error) {
	rows, err := s.db.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

type migrationFile struct {
	Version int64
	Name    string
	Content string
}

func (s *MigrationService) getPendingMigrations(applied map[int64]bool) ([]migrationFile, error) {
	var pending []migrationFile

	files, err := fs.ReadDir(migrations.FS, "sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		// Parse version from filename (e.g., "001_initial_schema.up.sql")
		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 2 {
			continue // Skip invalid filenames
		}

		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue // Skip if version is not a number
		}

		if !strings.HasSuffix(name, ".up.sql") {
			continue // Only process up migrations
		}

		if applied[version] {
			continue
		}

		content, err := fs.ReadFile(migrations.FS, filepath.Join("sql", name))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		pending = append(pending, migrationFile{
			Version: version,
			Name:    name,
			Content: string(content),
		})
	}

	// Sort by version
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})

	return pending, nil
}

func (s *MigrationService) applyMigration(ctx context.Context, m migrationFile) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Execute migration SQL
	if _, err := tx.Exec(ctx, m.Content); err != nil {
		return fmt.Errorf("failed to execute migration sql: %w", err)
	}

	// Record version
	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, dirty, executed_at) VALUES ($1, $2, $3)", m.Version, false, time.Now()); err != nil {
		return fmt.Errorf("failed to record migration version: %w", err)
	}

	return tx.Commit(ctx)
}
