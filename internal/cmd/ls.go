package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/output"
	"github.com/keygen-sh/keygen-relay/internal/try"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"github.com/spf13/cobra"
)

func LsCmd(manager licenses.Manager) *cobra.Command {
	var (
		plain bool
		pool  *string
	)

	cmd := &cobra.Command{
		Use:          "ls",
		Short:        "print the local relay server's license pool, with stats for each license",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if p, err := cmd.Flags().GetString("pool"); err == nil {
				if p != "" {
					pool = &p
				}
			}

			poolList, err := manager.GetPools(cmd.Context())
			if err != nil {
				return err
			}

			pools := make(map[int64]string, len(poolList))
			for _, p := range poolList {
				pools[p.ID] = p.Name
			}

			licensesList, err := manager.ListLicenses(cmd.Context(), pool)
			if err != nil {
				output.PrintError(cmd.ErrOrStderr(), err.Error())

				return nil
			}

			if len(licensesList) == 0 {
				output.PrintSuccess(cmd.OutOrStdout(), "license pool is empty")

				return nil
			}

			var renderer ui.TableRenderer
			if plain {
				renderer = ui.NewSimpleTableRenderer(cmd.OutOrStdout())
			} else {
				renderer = ui.NewBubbleteaTableRenderer()
			}

			columns := []table.Column{
				{Title: "id", Width: 36},
				{Title: "pool", Width: 8}, // start with min width
				{Title: "claims", Width: 8},
				{Title: "node_id", Width: 8},
				{Title: "last_claimed_at", Width: 20},
				{Title: "last_released_at", Width: 20},
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

				var poolStr string
				if lic.PoolID != nil {
					if name, ok := pools[*lic.PoolID]; ok {
						poolStr = name
					} else {
						poolStr = "<n/a>" // should never happen
					}
				} else {
					poolStr = "-"
				}

				// update pool column width dynamically
				if poolWidth := len(poolStr); poolWidth > columns[1].Width && poolWidth <= 32 {
					columns[1].Width = poolWidth
				} else if poolWidth > 32 {
					columns[1].Width = 32
				}

				lastClaimedAtStr := formatTime(lic.LastClaimedAt)
				lastReleasedAtStr := formatTime(lic.LastReleasedAt)

				tableRows = append(tableRows, table.Row{lic.Guid, poolStr, claimsStr, nodeIDStr, lastClaimedAtStr, lastReleasedAtStr})
			}

			if err := renderer.Render(tableRows, columns); err != nil {
				output.PrintError(cmd.ErrOrStderr(), fmt.Sprintf("error rendering table: %v", err))

				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&plain, "plain", try.Try(try.EnvBool("RELAY_PLAIN"), try.EnvBool("NO_COLOR"), try.EnvBool("CI"), try.Static(false)), "display the table in plain text format [$RELAY_PLAIN=1]")
	cmd.Flags().String("pool", try.Try(try.Env("RELAY_POOL"), try.Static("")), "pool to list licenses from [$RELAY_POOL=prod]")

	_ = cmd.RegisterFlagCompletionFunc("pool", poolTypeCompletion)

	return cmd
}

func formatTime(t *int64) string {
	if t == nil {
		return "-"
	}

	return time.Unix(*t, 0).UTC().Format(time.RFC3339)
}
