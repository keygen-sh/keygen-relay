package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "version",
		Short:        "print the current relay version",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Print(cmd.OutOrStdout(), cmd.Root().Version)

			return nil
		},
	}

	return cmd
}
