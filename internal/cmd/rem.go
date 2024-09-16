package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func RemCmd(manager licenses.Manager) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "rem",
		Short: "Delete a license from the local relay server's pool",
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.RemoveLicense(cmd.Context(), id)
			if err != nil {
				fmt.Println("error creating license record", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "license ID to remove")
	cmd.MarkFlagRequired("id")

	return cmd
}
