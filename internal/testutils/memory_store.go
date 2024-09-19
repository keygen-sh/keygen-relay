package testutils

import (
	"database/sql"
	"fmt"
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

	// Print to debug the opening of a new connection
	fmt.Println("Opening new in-memory database connection")

	schema, err := os.ReadFile("../../db/schema.sql")
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}
	require.NoError(t, err)

	// Print to debug schema application
	fmt.Println("Applying schema")
	if _, err := dbConn.Exec(string(schema)); err != nil {
		t.Fatalf("Failed to apply schema: %v", err)
	}

	store := db.NewStore(db.New(dbConn), dbConn)

	return store, dbConn
}

func CloseMemoryStore(dbConn *sql.DB) {
	// Print to debug the closing of the connection
	fmt.Println("Closing in-memory database connection")
	if err := dbConn.Close(); err != nil {
		log.Printf("Failed to close in-memory database connection: %v", err)
	}
}
