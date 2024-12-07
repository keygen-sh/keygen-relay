package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/spf13/cobra"
)

func DelCmd(manager licenses.Manager) *cobra.Command {
	var licenseID string

	cmd := &cobra.Command{
		Use:          "del",
		Short:        "delete a license from the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := manager.RemoveLicense(cmd.Context(), licenseID)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())
				return nil
			}

			output.PrintSuccess(cmd.OutOrStdout(), "license deleted successfully")

			return nil
		},
	}

	cmd.Flags().StringVar(&licenseID, "license", "", "license ID to remove")
	_ = cmd.MarkFlagRequired("license")

	return cmd
}
