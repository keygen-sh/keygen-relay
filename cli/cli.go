package cli

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"

	schema "github.com/keygen-sh/keygen-relay/db"
	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	// The current version of Relay, embedded at compile time.
	Version = "<not set>"
)

func Run() int {
	var dbConnection *sql.DB

	cfg := config.New()
	manager := licenses.NewManager(cfg.License, os.ReadFile, licenses.NewKeygenLicenseVerifier)
	srv := server.New(cfg.Server, manager)

	rootCmd := &cobra.Command{
		Use:   "relay",
		Short: "Keygen Relay CLI",
		Long: `relay is a small command line utility that distributes license files to nodes on a local network

Version:
  relay/` + Version + " " + runtime.GOOS + "-" + runtime.GOARCH + " " + runtime.Version(),
		SilenceUsage: true,
		Version:      Version,
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
	rootCmd.PersistentFlags().CountVarP(&cfg.Logger.Verbosity, "verbose", "v", `log verbosity e.g. -vvv (default "error")`)
	rootCmd.PersistentFlags().Bool("no-audit", false, "disable audit logs")
	rootCmd.PersistentFlags().BoolVar(&cfg.Logger.DisableColor, "no-color", false, "disable colors in command output [$NO_COLOR=1]")

	rootCmd.AddCommand(cmd.AddCmd(manager))
	rootCmd.AddCommand(cmd.DelCmd(manager))
	rootCmd.AddCommand(cmd.LsCmd(manager))
	rootCmd.AddCommand(cmd.StatCmd(manager))
	rootCmd.AddCommand(cmd.ServeCmd(srv))
	rootCmd.AddCommand(cmd.VersionCmd())

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

	slog.Info("applying database pragmas", "path", cfg.DB.DatabaseFilePath)

	// set the journal mode to Write-Ahead Logging for concurrency
	if _, err := dbConn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Fatalf("failed to set journal_mode pragma: %v", err)
	}

	// set synchronous mode to NORMAL to better balance performance and safety
	if _, err := dbConn.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		log.Fatalf("failed to set synchronous pragma: %v", err)
	}

	// set busy timeout to 5 seconds to avoid lock-related errors
	if _, err := dbConn.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		log.Fatalf("failed to set busy_timeout pragma: %v", err)
	}

	// set cache size to 20MB for faster data access
	if _, err := dbConn.Exec("PRAGMA cache_size = -20000"); err != nil {
		log.Fatalf("failed to set cache_size pragma: %v", err)
	}

	// enable foreign key constraints
	if _, err := dbConn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("failed to set foreign_keys pragma: %v", err)
	}

	// enable auto vacuuming and set it to incremental mode for gradual space reclaiming
	if _, err := dbConn.Exec("PRAGMA auto_vacuum = INCREMENTAL"); err != nil {
		log.Fatalf("failed to set auto_vacuum pragma: %v", err)
	}

	// store temporary tables and data in memory for better performance
	if _, err := dbConn.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		log.Fatalf("failed to set temp_store pragma: %v", err)
	}

	// set the mmap_size to 2GB for faster read/write access using memory-mapped I/O
	if _, err := dbConn.Exec("PRAGMA mmap_size = 2147483648"); err != nil {
		log.Fatalf("failed to set mmap_size pragma: %v", err)
	}

	// set the page size to 8KB for balanced memory usage and performance
	if _, err := dbConn.Exec("PRAGMA page_size = 8192"); err != nil {
		log.Fatalf("failed to set page_size pragma: %v", err)
	}

	if !dbExists {
		slog.Info("applying database schema", "path", cfg.DB.DatabaseFilePath)

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
