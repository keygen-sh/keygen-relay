package cli

import (
	"context"
	"database/sql"
	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var dbConnection *sql.DB

func Run() int {
	cfg := config.New()
	ctx := context.Background()

	manager := licenses.NewManager(cfg.License, os.ReadFile, licenses.NewKeygenLicenseVerifier)

	rootCmd := &cobra.Command{
		Use:          "relay",
		Short:        "Keygen Relay CLI",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger.Init(cfg.Logger, os.Stdout)

			// Initialization database connection in PersistentPreRun hook for getting persistent flags
			store, dbConn, err := initStore(ctx, cfg)
			if err != nil {
				slog.Error("failed to initialize store", "error", err)
				return err
			}

			dbConnection = dbConn

			manager.AttachStore(store)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if dbConnection != nil {
				if err := dbConnection.Close(); err != nil {
					slog.Error("failed to close database connection", "error", err)
					return err
				}
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfg.DB.DatabaseFilePath, "database", "./relay.sqlite", "specify an alternate database path")
	rootCmd.PersistentFlags().CountVarP(&cfg.Logger.Verbosity, "verbose", "v", "counted verbosity")

	tableRenderer := ui.NewBubbleteaTableRenderer()

	rootCmd.AddCommand(cmd.AddCmd(manager))
	rootCmd.AddCommand(cmd.DelCmd(manager))
	rootCmd.AddCommand(cmd.LsCmd(manager, tableRenderer))
	rootCmd.AddCommand(cmd.StatCmd(manager, tableRenderer))

	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute command", "error", err)
		return 1
	}

	return 0
}

func initStore(ctx context.Context, cfg *config.Config) (licenses.Store, *sql.DB, error) {
	dbExists := fileExists(cfg.DB.DatabaseFilePath)
	dbConn, err := sql.Open("sqlite3", cfg.DB.DatabaseFilePath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, nil, err
	}

	if err := dbConn.Ping(); err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, nil, err
	}

	if !dbExists {
		slog.Info("database does not exist, applying schema")
		schema, err := os.ReadFile("db/schema.sql")
		if err != nil {
			slog.Error("failed to read schema file", "error", err)
			return nil, nil, err
		}

		if _, err := dbConn.ExecContext(ctx, string(schema)); err != nil {
			slog.Error("failed to execute schema", "error", err)
			return nil, nil, err
		}
	}

	queries := db.New(dbConn)
	return db.NewStore(queries, dbConn), dbConn, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		slog.Warn("file does not exist", "filename", filename)
		return false
	}
	return !info.IsDir()
}
