package cmd

import (
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/locker"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/spf13/cobra"
)

func AddCmd(manager licenses.Manager) *cobra.Command {
	var (
		publicKey = locker.PublicKey
		filePath  string
		key       string
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

			output.PrintSuccess(cmd.OutOrStdout(), "license added successfully")

			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "path to a signed and encrypted license file")
	cmd.Flags().StringVar(&key, "key", "", "license key for decryption")

	if !locker.Locked() {
		cmd.Flags().StringVar(&publicKey, "public-key", try.Try(try.Env("RELAY_PUBLIC_KEY"), try.Static("")), "your keygen.sh public key for verification [$KEYGEN_PUBLIC_KEY=e8601...e48b6]")
	}

	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("public-key")

	return cmd
}
