package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/spf13/cobra"
)

func AddCmd(manager licenses.Manager) *cobra.Command {
	var (
		filePath  string
		key       string
		publicKey string
	)

	cmd := &cobra.Command{
		Use:          "add",
		Short:        "push a license to the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := manager.AddLicense(cmd.Context(), filePath, key, publicKey)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			output.PrintSuccess(cmd.OutOrStdout(), "License added successfully.")

			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "license file")
	cmd.Flags().StringVar(&key, "key", "", "license key")
	cmd.Flags().StringVar(&publicKey, "public-key", "", "public key for cryptographically verified")

	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("public-key")

	return cmd
}
