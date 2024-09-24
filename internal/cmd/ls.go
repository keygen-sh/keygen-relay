package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func LsCmd(manager licenses.Manager, tableRenderer ui.TableRenderer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ls",
		Short:        "Print the local relay server's license pool, with stats for each license",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			licensesList, err := manager.ListLicenses(cmd.Context())
			if err != nil {
				return err
			}

			if len(licensesList) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No licenses found.")

				return nil
			}

			columns := []table.Column{
				{Title: "ID", Width: 36},
				{Title: "Claims", Width: 8},
				{Title: "NodeID", Width: 8},
				{Title: "Last Claimed At", Width: 20},
				{Title: "Last Released At", Width: 20},
			}

			tableRows := make([]table.Row, 0, len(licensesList))
			for _, lic := range licensesList {
				claimsStr := fmt.Sprintf("%d", lic.Claims)

				var nodeIDStr string
				if lic.NodeID != nil {
					nodeIDStr = strconv.FormatInt(*lic.NodeID, 10)
				} else {
					nodeIDStr = "-"
				}

				lastClaimedAtStr := formatTime(lic.LastClaimedAt)
				lastReleasedAtStr := formatTime(lic.LastReleasedAt)

				tableRows = append(tableRows, table.Row{lic.ID, claimsStr, nodeIDStr, lastClaimedAtStr, lastReleasedAtStr})
			}

			if err := tableRenderer.Render(tableRows, columns); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error rendering table: %v", err)
				return err
			}

			return nil
		},
	}

	return cmd
}

func formatTime(t *string) string {
	if t == nil {
		return "-"
	}

	parsedTime, err := time.Parse(time.RFC3339, *t)
	if err != nil {
		return "-"
	}

	return parsedTime.Format("2006-01-02 15:04:05")
}
