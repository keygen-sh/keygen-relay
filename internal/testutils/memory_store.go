package testutils

import (
	"database/sql"
	schema "github.com/keygen-sh/keygen-relay/db"
	"log"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func NewMemoryStore(t *testing.T) (*db.Store, *sql.DB) {
	dbConn, err := sql.Open("sqlite3", ":memory:?_pragma=foreign_keys(on)")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	_, err = dbConn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	if _, err := dbConn.Exec(schema.SchemaSQL); err != nil {
		t.Fatalf("Failed to apply schema: %v", err)
	}

	store := db.NewStore(db.New(dbConn), dbConn)

	return store, dbConn
}

func CloseMemoryStore(dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		log.Printf("Failed to close in-memory database connection: %v", err)
	}
}
