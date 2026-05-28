package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Migrate applies all pending "*.up.sql" migrations in lexical order. Applied
// versions are tracked in the schema_migrations table so the runner is
// idempotent and safe to call on every startup.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var versions []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			versions = append(versions, e.Name())
		}
	}
	sort.Strings(versions)

	for _, file := range versions {
		var exists bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`,
			file).Scan(&exists); err != nil {
			return fmt.Errorf("check migration %s: %w", file, err)
		}
		if exists {
			continue
		}

		sqlBytes, err := migrationFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply migration %s: %w", file, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations(version) VALUES ($1)`, file); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", file, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", file, err)
		}
	}
	return nil
}
