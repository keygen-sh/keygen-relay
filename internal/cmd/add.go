package cmd

import (
	"strings"

	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/locker"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/spf13/cobra"
)

func AddCmd(manager licenses.Manager) *cobra.Command {
	var (
		publicKey = locker.PublicKey
	)

	cmd := &cobra.Command{
		Use:          "add",
		Short:        "push license(s) to the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			publicKey = strings.TrimSpace(publicKey)

			files, err := cmd.Flags().GetStringSlice("file")
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			keys, err := cmd.Flags().GetStringSlice("key")
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			if len(files) != len(keys) {
				output.PrintError(cmd.ErrOrStderr(), "number of key and file flags must match")

				return nil
			}

			for i := range len(files) {
				file := files[i]
				key := strings.TrimSpace(keys[i])

				license, err := manager.AddLicense(cmd.Context(), file, key, publicKey)
				if err != nil {
					output.PrintError(cmd.ErrOrStderr(), err.Error())

					return nil
				}

				output.PrintSuccess(cmd.OutOrStdout(), "license added successfully: %s", license.Guid)
			}

			return nil
		},
	}

	cmd.Flags().StringSlice("file", nil, "path to a signed and encrypted license file")
	cmd.Flags().StringSlice("key", nil, "license key for decryption")

	if !locker.Locked() {
		cmd.Flags().StringVar(&publicKey, "public-key", try.Try(try.Env("RELAY_PUBLIC_KEY"), try.Static("")), "your keygen.sh public key for verification [$KEYGEN_PUBLIC_KEY=e860..48b6]")
	}

	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("public-key")

	return cmd
}
