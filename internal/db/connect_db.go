package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"log"
)

func ConnectDB(dbPath string) (*sql.DB, error) {
	ctx := context.Background()

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := dbConn.PingContext(ctx); err != nil {
		return nil, err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	n, err := migrate.Exec(dbConn, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return dbConn, err
	}

	fmt.Printf("Applied %d migrations!\n", n)

	log.Printf("Connected to database: %s", dbPath)

	return dbConn, nil
}
