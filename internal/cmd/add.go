package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/spf13/cobra"
	"os"
)

func addCmd(queries *db.Queries) *cobra.Command {
	var (
		filePath  string
		key       string
		publicKey string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Push a license to the local relay server's pool",
		Run: func(cmd *cobra.Command, args []string) {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}

			//TODO: add verifying, decrypting and move it to service
			id := uuid.New().String()
			err = queries.InsertLicense(cmd.Context(), db.InsertLicenseParams{File: fileContent, Key: key, ID: id})
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
