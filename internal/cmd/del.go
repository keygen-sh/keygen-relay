package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func DelCmd(manager licenses.Manager) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "del",
		Short: "Delete a license from the local relay server's pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := manager.RemoveLicense(cmd.Context(), id)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error removing license record: %v\n", err)
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "License removed successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "license ID to remove")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
