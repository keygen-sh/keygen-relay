package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/spf13/cobra"
)

func strategyTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"fifo", "lifo", "rand"}, cobra.ShellCompDirectiveDefault
}

func poolTypeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	pools, err := getPoolNamesForCompletion(cmd)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	return pools, cobra.ShellCompDirectiveDefault
}

func getPoolNamesForCompletion(cmd *cobra.Command) ([]string, error) {
	ctx := cmd.Context()
	path := try.Try(
		try.CmdPersistentFlag(cmd, "database"),
		try.Env("RELAY_DATABASE"),
		try.Static("./relay.sqlite"),
	)

	store, conn, err := db.Connect(ctx, &db.Config{DatabaseFilePath: path})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pools, err := store.GetPools(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(pools))
	for i, pool := range pools {
		names[i] = pool.Name
	}

	return names, nil
}
