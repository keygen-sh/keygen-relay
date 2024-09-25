package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"github.com/spf13/cobra"
	"strconv"
)

func StatCmd(manager licenses.Manager) *cobra.Command {
	var licenseID string
	var plain bool

	cmd := &cobra.Command{
		Use:          "stat",
		Short:        "Print stats for a license in the local relay server's pool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			license, err := manager.GetLicenseByID(cmd.Context(), licenseID)
			if err != nil {
				return err
			}

			columns := []table.Column{
				{Title: "ID", Width: 36},
				{Title: "Claims", Width: 8},
				{Title: "NodeID", Width: 8},
				{Title: "Last Claimed At", Width: 20},
				{Title: "Last Released At", Width: 20},
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
				fmt.Fprintf(cmd.ErrOrStderr(), "Error rendering table: %v", err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&licenseID, "id", "", "License ID to print stats for")
	_ = cmd.MarkFlagRequired("id")

	cmd.Flags().BoolVar(&plain, "plain", false, "display the table in plain text format")

	return cmd
}
