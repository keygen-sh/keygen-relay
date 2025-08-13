package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/spf13/cobra"
)

func DelCmd(manager licenses.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "del",
		Short:        "delete license(s) from the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			licenseIDs, err := cmd.Flags().GetStringSlice("license")
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			for _, licenseID := range licenseIDs {
				err := manager.RemoveLicense(cmd.Context(), licenseID)
				if err != nil {
					output.PrintError(cmd.ErrOrStderr(), err.Error())

					return nil
				}

				output.PrintSuccess(cmd.OutOrStdout(), "license deleted successfully: %s", licenseID)
			}

			return nil
		},
	}

	cmd.Flags().StringSlice("license", nil, "license ID to remove")
	_ = cmd.MarkFlagRequired("license")

	return cmd
}
