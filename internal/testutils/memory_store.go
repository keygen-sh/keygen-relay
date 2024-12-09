package testutils

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	schema "github.com/keygen-sh/keygen-relay/db"
	"github.com/keygen-sh/keygen-relay/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func NewMemoryStore(t *testing.T) (*db.Store, *sql.DB) {
	conn, err := sql.Open("sqlite3", ":memory:?_pragma=foreign_keys(on)")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	_, err = conn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	migrations, err := iofs.New(schema.Migrations, "migrations")
	if err != nil {
		t.Fatalf("failed to initialize migrations fs: %v", err)
	}

	migrator, err := db.NewMigrator(conn, migrations)
	if err != nil {
		t.Fatalf("failed to initialize migrations: %v", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	store := db.NewStore(db.New(conn), conn)

	return store, conn
}

func CloseMemoryStore(dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		log.Printf("failed to close in-memory database connection: %v", err)
	}
}
