package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	schema "github.com/keygen-sh/keygen-relay/db"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(ctx context.Context, cfg *Config) (*Store, *sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_txlock=immediate", cfg.DatabaseFilePath)
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()

		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// apply database pragmas e.g. WAL mode, timeouts, etc.
	logger.Info("applying database pragmas", "path", cfg.DatabaseFilePath)

	for key, value := range cfg.DatabasePragmas {
		if _, err := conn.Exec(fmt.Sprintf("PRAGMA %s = %s", key, value)); err != nil {
			conn.Close()

			return nil, nil, fmt.Errorf("failed to set pragma %s=%s: %w", key, value, err)
		}
	}

	// apply migrations e.g. initial schema, new schema changes, etc.
	logger.Info("applying database migrations", "path", cfg.DatabaseFilePath)

	migrations, err := iofs.New(schema.Migrations, "migrations")
	if err != nil {
		conn.Close()

		return nil, nil, fmt.Errorf("failed to initialize migrations fs: %w", err)
	}

	migrator, err := NewMigrator(conn, migrations)
	if err != nil {
		conn.Close()

		return nil, nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		conn.Close()

		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	queries := New(conn)
	store := NewStore(queries, conn)

	return store, conn, nil
}
