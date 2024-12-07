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
	var licenseID string
	var plain bool

	cmd := &cobra.Command{
		Use:          "stat",
		Short:        "print stats for a license in the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			license, err := manager.GetLicenseByID(cmd.Context(), licenseID)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			columns := []table.Column{
				{Title: "id", Width: 36},
				{Title: "claims", Width: 8},
				{Title: "node_id", Width: 8},
				{Title: "last_claimed_at", Width: 20},
				{Title: "last_released_at", Width: 20},
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
				{license.ID, claimsStr, nodeIDStr, lastClaimedAtStr, lastReleasedAtStr},
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
