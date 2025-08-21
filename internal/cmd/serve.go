package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/locker"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/spf13/cobra"
)

const minTTL = 30 * time.Second

func ServeCmd(srv server.Server) *cobra.Command {
	cfg := srv.Config()

	handler := server.NewHandler(srv)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()

		logger.Debug("route registered", "path", path)

		return nil
	})

	router.Use(server.LoggingMiddleware)

	// Mount the router to the server
	srv.Mount(router)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "run the relay server to manage license distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			if ttl, err := cmd.Flags().GetDuration("ttl"); err == nil {
				if err := validateTTL(ttl); err != nil {
					output.PrintError(cmd.ErrOrStderr(), err.Error())

					return err
				}
			}

			if disableHeartbeats, err := cmd.Flags().GetBool("no-heartbeats"); err == nil {
				cfg.EnabledHeartbeat = !disableHeartbeats
			}

			// workaround for lack of support for nullable string flags
			if p, err := cmd.Flags().GetString("pool"); err == nil {
				if p != "" {
					cfg.Pool = &p
				}
			}

			srv.Manager().Config().Strategy = string(cfg.Strategy)
			srv.Manager().Config().ExtendOnHeartbeat = cfg.EnabledHeartbeat

			output.PrintSuccess(cmd.OutOrStdout(), "the server is starting")

			if err := srv.Run(); err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			return nil
		},
	}

	// FIXME(ezekg) add default strategy since Var() doesn't support a default
	cfg.Strategy = try.Try(
		try.EnvAs("RELAY_STRATEGY", func(value string) server.StrategyType {
			return server.StrategyType(value)
		}),
		try.Static(cfg.Strategy),
	)

	if locker.LockedAddr() {
		cfg.ServerAddr = locker.Addr
	} else {
		cmd.Flags().StringVarP(&cfg.ServerAddr, "bind", "b", try.Try(try.Env("RELAY_ADDR"), try.Env("BIND_ADDR"), try.Static(cfg.ServerAddr)), "ip address to bind the relay server to [$RELAY_ADDR=0.0.0.0]")
	}

	if locker.LockedPort() {
		port, err := strconv.Atoi(locker.Port)
		if err != nil {
			panic(err)
		}

		cfg.ServerPort = port
	} else {
		cmd.Flags().IntVarP(&cfg.ServerPort, "port", "p", try.Try(try.EnvInt("RELAY_PORT"), try.EnvInt("PORT"), try.Static(cfg.ServerPort)), "port to run the relay server on [$RELAY_PORT=6349]")
	}

	cmd.Flags().DurationVar(&cfg.TTL, "ttl", try.Try(try.EnvDuration("RELAY_LEASE_TTL"), try.Static(cfg.TTL)), "time-to-live for leases [$RELAY_LEASE_TTL=60s]")
	cmd.Flags().Bool("no-heartbeats", try.Try(try.EnvBool("RELAY_NO_HEARTBEATS"), try.Static(false)), "disable node heartbeat monitoring and culling as well as lease extensions [$RELAY_NO_HEARTBEAT=1]")
	cmd.Flags().Var(&cfg.Strategy, "strategy", `strategy for license distribution e.g. "fifo", "lifo", or "rand" [$RELAY_STRATEGY=rand]`)
	cmd.Flags().DurationVar(&cfg.CullInterval, "cull-interval", try.Try(try.EnvDuration("RELAY_CULL_INTERVAL"), try.Static(cfg.CullInterval)), "interval at which to cull dead nodes [$RELAY_CULL_INTERVAL=15s]")
	cmd.Flags().String("pool", try.Try(try.Env("RELAY_POOL"), try.Static("")), "pool to serve licenses from [$RELAY_POOL=prod]")

	_ = cmd.RegisterFlagCompletionFunc("strategy", strategyTypeCompletion)
	_ = cmd.RegisterFlagCompletionFunc("pool", poolTypeCompletion)

	return cmd
}

func validateTTL(ttl time.Duration) error {
	if ttl < minTTL {
		return fmt.Errorf("time-to-live value must be at least %s", minTTL)
	}
	return nil
}

func strategyTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"fifo", "lifo", "rand"}, cobra.ShellCompDirectiveDefault
}

func poolTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// TODO(ezekg) query pools
	return []string{}, cobra.ShellCompDirectiveDefault
}
