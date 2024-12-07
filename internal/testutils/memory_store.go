package testutils

import (
	"database/sql"
	"log"
	"testing"

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

	if _, err := conn.Exec(schema.Schema); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}

	store := db.NewStore(db.New(conn), conn)

	return store, conn
}

func CloseMemoryStore(dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		log.Printf("failed to close in-memory database connection: %v", err)
	}
}
