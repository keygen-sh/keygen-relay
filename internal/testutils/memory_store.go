package testutils

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func NewMemoryStore(t *testing.T) (*db.Store, *sql.DB) {
	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	schema, err := os.ReadFile("../../db/schema.sql")
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}
	require.NoError(t, err)

	if _, err := dbConn.Exec(string(schema)); err != nil {
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
