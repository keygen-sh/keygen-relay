package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/spf13/cobra"
)

const minTTL = 30 * time.Second

func ServeCmd(srv server.Server) *cobra.Command {
	cfg := srv.Config()

	handler := server.NewHandler(srv.Manager())
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		slog.Debug("Route registered", "path", path)
		return nil
	})

	router.Use(server.LoggingMiddleware)

	// Mount the router to the server
	srv.Mount(router)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "run the relay server to manage license distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			disableHeartbeat, err := cmd.Flags().GetBool("no-heartbeats")
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), fmt.Sprintf("Failed to parse 'no-heartbeats' flag: %v", err))
				return err
			}

			cfg.EnabledHeartbeat = !disableHeartbeat

			ttl, err := cmd.Flags().GetDuration("ttl")
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), fmt.Sprintf("Failed to parse 'ttl' flag: %v", err))
				return err
			}

			if err := validateTTL(ttl); err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())
				return err
			}

			srv.Manager().Config().Strategy = string(cfg.Strategy)
			srv.Manager().Config().ExtendOnHeartbeat = cfg.EnabledHeartbeat

			output.PrintSuccess(cmd.OutOrStdout(), "The server is starting")

			if err := srv.Run(); err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())
				return nil
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&cfg.ServerPort, "port", "p", cfg.ServerPort, "port to run the relay server on")
	cmd.Flags().DurationVarP(&cfg.TTL, "ttl", "t", cfg.TTL, "time-to-live for license claims")
	cmd.Flags().Bool("no-heartbeats", false, "disable heartbeat mechanism")
	cmd.Flags().Var(&cfg.Strategy, "strategy", `strategy type for license distribution e.g. "fifo", "lifo", or "rand"`)
	cmd.Flags().DurationVar(&cfg.CleanupInterval, "cleanup-interval", cfg.CleanupInterval, "interval at which to cull inactive nodes.")

	_ = cmd.RegisterFlagCompletionFunc("strategy", strategyTypeCompletion)

	return cmd
}

func validateTTL(ttl time.Duration) error {
	if ttl < minTTL {
		return fmt.Errorf("TTL value must be at least %s", minTTL)
	}
	return nil
}

func strategyTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"fifo", "lifo", "rand"}, cobra.ShellCompDirectiveDefault
}
