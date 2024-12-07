package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func (m Migrator) Up() error {
	return m.migrate.Up()
}

func (m Migrator) Down() error {
	return m.migrate.Down()
}

func NewMigrator(db *sql.DB, migrations source.Driver) (*Migrator, error) {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"file",
		migrations,
		"sqlite3",
		instance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrate: m}, nil
}
