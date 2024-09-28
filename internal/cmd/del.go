package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/spf13/cobra"
)

func DelCmd(manager licenses.Manager) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:          "del",
		Short:        "Delete a license from the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := manager.RemoveLicense(cmd.Context(), id)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())
				return nil
			}

			output.PrintSuccess(cmd.OutOrStdout(), "License deleted successfully.")

			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "license ID to remove")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
