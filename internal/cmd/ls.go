package cmd

import (
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func LsCmd(manager licenses.Manager, tableRenderer ui.TableRenderer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Print the local relay server's license pool, with stats for each license",
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
				{Title: "Key", Width: 50},
				{Title: "Claims", Width: 6},
				{Title: "NodeID", Width: 8},
			}

			tableRows := make([]table.Row, 0, len(licensesList))
			for _, lic := range licensesList {
				claimsStr := fmt.Sprintf("%d", lic.Claims)

				var nodeIDStr string
				if lic.NodeID.Valid {
					nodeIDStr = strconv.FormatInt(lic.NodeID.Int64, 10)
				} else {
					nodeIDStr = "-"
				}

				tableRows = append(tableRows, table.Row{lic.ID, lic.Key, claimsStr, nodeIDStr})
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
