package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"github.com/spf13/cobra"
)

func StatCmd(manager licenses.Manager) *cobra.Command {
	var (
		licenseID string
		plain     bool
	)

	cmd := &cobra.Command{
		Use:          "stat",
		Short:        "print stats for a license in the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			license, err := manager.GetLicenseByGUID(cmd.Context(), nil, licenseID)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			columns := []table.Column{
				{Title: "id", Width: 36},
				{Title: "pool", Width: 8}, // start with min width
				{Title: "claims", Width: 8},
				{Title: "node_id", Width: 8},
				{Title: "last_claimed_at", Width: 20},
				{Title: "last_released_at", Width: 20},
			}

			var poolStr string
			if license.PoolID != nil {
				pool, err := manager.GetPoolByID(cmd.Context(), *license.PoolID)
				if err != nil {
					return err
				}

				poolStr = pool.Name
			} else {
				poolStr = "-"
			}

			// update pool column width dynamically
			if poolWidth := len(poolStr); poolWidth > columns[1].Width && poolWidth <= 32 {
				columns[1].Width = poolWidth
			} else if poolWidth > 32 {
				columns[1].Width = 32
			}

			claimsStr := fmt.Sprintf("%d", license.Claims)

			var nodeIDStr string
			if license.NodeID != nil {
				nodeIDStr = strconv.FormatInt(*license.NodeID, 10)
			} else {
				nodeIDStr = "-"
			}

			lastClaimedAtStr := formatTime(license.LastClaimedAt)
			lastReleasedAtStr := formatTime(license.LastReleasedAt)

			tableRows := []table.Row{
				{license.Guid, poolStr, claimsStr, nodeIDStr, lastClaimedAtStr, lastReleasedAtStr},
			}

			var renderer ui.TableRenderer
			if plain {
				renderer = ui.NewSimpleTableRenderer(cmd.OutOrStdout())
			} else {
				renderer = ui.NewBubbleteaTableRenderer()
			}

			if err := renderer.Render(tableRows, columns); err != nil {
				output.PrintError(cmd.ErrOrStderr(), fmt.Sprintf("error rendering table: %v", err))

				return nil
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&licenseID, "license", "", "license ID to print stats for")
	_ = cmd.MarkFlagRequired("license")

	cmd.Flags().BoolVar(&plain, "plain", try.Try(try.EnvBool("RELAY_PLAIN"), try.EnvBool("NO_COLOR"), try.EnvBool("CI"), try.Static(false)), "display the table in plain text format [$RELAY_PLAIN=1]")

	return cmd
}
