package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/spf13/cobra"
	"log/slog"
)

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

	// Mount the router to the server
	srv.Mount(router)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the relay server to manage license distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			disableHeartbeat, err := cmd.Flags().GetBool("no-heartbeats")
			if err != nil {
				return fmt.Errorf("failed to parse 'no-heartbeats' flag: %v", err)
			}

			cfg.EnabledHeartbeat = !disableHeartbeat

			srv.Manager().Config().Strategy = string(cfg.Strategy)
			srv.Manager().Config().ExtendOnHeartbeat = cfg.EnabledHeartbeat

			if err := srv.Run(); err != nil {
				return fmt.Errorf("error running server: %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "The server is starting")
			return nil
		},
	}

	cmd.Flags().IntVarP(&cfg.ServerPort, "port", "p", cfg.ServerPort, "Port to run the relay server on")
	cmd.Flags().DurationVarP(&cfg.TTL, "ttl", "t", cfg.TTL, "Time-to-live for license claims")
	cmd.Flags().Bool("no-heartbeats", false, "Disable heartbeat mechanism")
	cmd.Flags().Var(&cfg.Strategy, "strategy", `Strategy type for license distribution. Allowed: "fifo", "lifo", "rand"`)
	cmd.Flags().DurationVar(&cfg.CleanupInterval, "cleanup-interval", cfg.CleanupInterval, "interval at which to check for inactive nodes.")

	_ = cmd.RegisterFlagCompletionFunc("strategy", strategyTypeCompletion)

	return cmd
}

func strategyTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"fifo", "lifo", "rand"}, cobra.ShellCompDirectiveDefault
}
