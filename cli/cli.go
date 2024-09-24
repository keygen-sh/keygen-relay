package cli

import (
	"context"
	"database/sql"
	"fmt"
	schema "github.com/keygen-sh/keygen-relay/db"
	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
)

func Run() int {
	var dbConnection *sql.DB

	cfg := config.New()
	manager := licenses.NewManager(cfg.License, os.ReadFile, licenses.NewKeygenLicenseVerifier)
	srv := server.New(cfg.Server, manager)

	rootCmd := &cobra.Command{
		Use:          "relay",
		Short:        "Keygen Relay CLI",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			logger.Init(cfg.Logger, os.Stdout)

			disableAudit, err := cmd.Flags().GetBool("no-audit")
			if err != nil {
				return fmt.Errorf("failed to parse 'no-audit' flag: %v", err)
			}
			cfg.License.EnabledAudit = !disableAudit

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
	rootCmd.PersistentFlags().Bool("no-audit", false, "disable audit logs")
	rootCmd.PersistentFlags().BoolVar(&cfg.Logger.DisableColor, "no-color", false, "Disable color logs")

	rootCmd.AddCommand(cmd.AddCmd(manager))
	rootCmd.AddCommand(cmd.DelCmd(manager))
	rootCmd.AddCommand(cmd.LsCmd(manager))
	rootCmd.AddCommand(cmd.StatCmd(manager))
	rootCmd.AddCommand(cmd.ServeCmd(srv))

	if err := rootCmd.Execute(); err != nil {
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

	//enable foreign key for sqlite
	_, err = dbConn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	if !dbExists {
		slog.Info("Applying database schema", "path", cfg.DB.DatabaseFilePath)

		if _, err := dbConn.ExecContext(ctx, schema.SchemaSQL); err != nil {
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
		slog.Debug("file does not exist", "filename", filename)
		return false
	}
	return !info.IsDir()
}
