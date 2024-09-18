package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
	"os"
)

func DelCmd(manager licenses.Manager) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "del",
		Short: "Delete a license from the local relay server's pool",
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.RemoveLicense(cmd.Context(), id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating license record: %v", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "license ID to remove")
	cmd.MarkFlagRequired("id")

	return cmd
}
