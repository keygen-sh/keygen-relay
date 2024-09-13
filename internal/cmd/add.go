package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func AddCmd(manager licenses.Manager) *cobra.Command {
	var (
		filePath  string
		key       string
		publicKey string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Push a license to the local relay server's pool",
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.AddLicense(cmd.Context(), filePath, key, publicKey)
			if err != nil {
				fmt.Println("error creating license record", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "license file")
	cmd.Flags().StringVar(&key, "key", "", "license key")
	cmd.Flags().StringVar(&publicKey, "public-key", "", "public key for cryptographically verified")

	return cmd
}
