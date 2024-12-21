package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	schema "github.com/keygen-sh/keygen-relay/db"
	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/keygen-sh/keygen-relay/internal/try"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	// The current version of Relay, embedded at compile time.
	Version = "<not set>"
)

func Run() int {
	var conn *sql.DB

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

			if pragmas, err := cmd.Flags().GetStringSlice("pragma"); err == nil {
				for _, pragma := range pragmas {
					keyvalues := strings.SplitN(pragma, "=", 2)
					if len(keyvalues) != 2 {
						return fmt.Errorf("invalid pragma format: %s (expected key=value)", pragma)
					}

					key, value := keyvalues[0], keyvalues[1]
					if key == "" || value == "" {
						return fmt.Errorf("invalid pragma format: %s (expected key=value)", pragma)
					}

					cfg.DB.DatabasePragmas[key] = value
				}
			}

			if disableAudit, err := cmd.Flags().GetBool("no-audit"); err == nil {
				cfg.License.EnabledAudit = !disableAudit
			}

			// init database connection in PersistentPreRun hook for getting persistent flags
			var store *db.Store
			var err error

			store, conn, err = initStore(ctx, cfg)
			if err != nil {
				slog.Error("failed to initialize store", "error", err)

				return err
			}

			manager.AttachStore(*store)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if conn != nil {
				if err := conn.Close(); err != nil {
					slog.Error("failed to close database connection", "error", err)

					return err
				}
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfg.DB.DatabaseFilePath, "database", try.Try(try.Env("RELAY_DATABASE"), try.Static("./relay.sqlite")), "the path to a .sqlite database file [$RELAY_DATABASE=./relay.sqlite]")
	rootCmd.PersistentFlags().CountVarP(&cfg.Logger.Verbosity, "verbose", "v", `log level e.g. -vvv for "info" (default -v=1 i.e. "error") [$DEBUG=1]`)
	rootCmd.PersistentFlags().Bool("no-audit", try.Try(try.EnvBool("RELAY_NO_AUDIT"), try.Static(false)), "disable audit logs [$RELAY_NO_AUDIT=1]")
	rootCmd.PersistentFlags().BoolVar(&cfg.Logger.DisableColor, "no-color", false, "disable colors in command output [$NO_COLOR=1]")
	rootCmd.PersistentFlags().StringSlice("pragma", nil, "database pragma key-value pairs (e.g. --pragma mmap_size=536870912 --pragma synchronous=OFF)")

	rootCmd.SetHelpCommand(cmd.HelpCmd())
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

func initStore(_ context.Context, cfg *config.Config) (*db.Store, *sql.DB, error) {
	conn, err := sql.Open("sqlite3", cfg.DB.DatabaseFilePath)
	if err != nil {
		slog.Error("failed to open database", "error", err)

		return nil, nil, err
	}

	if err := conn.Ping(); err != nil {
		slog.Error("failed to connect to database", "error", err)

		return nil, nil, err
	}

	slog.Info("applying database pragmas", "path", cfg.DB.DatabaseFilePath)

	for key, value := range cfg.DB.DatabasePragmas {
		if _, err := conn.Exec(fmt.Sprintf("PRAGMA %s = %s", key, value)); err != nil {
			slog.Error("failed to set pragma", "key", key, "value", value, "error", err)

			return nil, nil, err
		}
	}

	// apply migrations e.g. initial schema, etc.
	slog.Info("applying database migrations", "path", cfg.DB.DatabaseFilePath)

	migrations, err := iofs.New(schema.Migrations, "migrations")
	if err != nil {
		slog.Error("failed to initialize migrations fs", "error", err)

		return nil, nil, err
	}

	migrator, err := db.NewMigrator(conn, migrations)
	if err != nil {
		slog.Error("failed to initialize migrations", "error", err)

		return nil, nil, err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("failed to apply migrations", "error", err)

		return nil, nil, err
	}

	queries := db.New(conn)
	store := db.NewStore(queries, conn)

	return store, conn, nil
}
