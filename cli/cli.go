package cli

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
)

func Run() int {
	cfg := config.New()
	ctx := context.Background()

	rootCmd := &cobra.Command{
		Use:   "relay",
		Short: "Keygen Relay CLI",
	}

	rootCmd.PersistentFlags().StringVar(&cfg.DB.DatabaseFilePath, "database", "./relay.sqlite", "specify an alternate database path (default: ./relay.sqlite)")
	rootCmd.PersistentFlags().CountVarP(&cfg.Logger.Verbosity, "verbose", "v", "counted verbosity")

	logger.Init(cfg.Logger.Verbosity)

	store, dbConn, err := initStore(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return 1
	}

	defer dbConn.Close()

	manager := licenses.NewManager(store, cfg.License)

	rootCmd.AddCommand(cmd.AddCmd(manager))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return 1
	}

	return 0
}

func initStore(ctx context.Context, cfg *config.Config) (licenses.Store, *sql.DB, error) {
	dbExists := fileExists(cfg.DB.DatabaseFilePath)
	dbConn, err := sql.Open("sqlite3", cfg.DB.DatabaseFilePath)
	if err != nil {
		return nil, nil, err
	}

	if err := dbConn.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		return nil, nil, err
	}

	if !dbExists {
		log.Println("Database does not exist, applying schema")

		schema, err := os.ReadFile("db/schema.sql")
		if err != nil {
			return nil, nil, err
		}

		if _, err := dbConn.ExecContext(ctx, string(schema)); err != nil {
			return nil, nil, err
		}
	}

	queries := db.New(dbConn)

	return db.NewStore(queries), dbConn, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
