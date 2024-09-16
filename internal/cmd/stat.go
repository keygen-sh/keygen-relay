package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"time"

	"database/sql"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keygen-sh/keygen-relay/internal/common"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/spf13/cobra"
)

func formatTime(t sql.NullString) string {
	if t.Valid {
		parsedTime, err := time.Parse(time.RFC3339, t.String)
		if err == nil {
			return parsedTime.Format("2000-01-01 00:00:00")
		}
	}
	return "-"
}

func StatCmd(manager licenses.Manager) *cobra.Command {
	var licenseID string

	cmd := &cobra.Command{
		Use:   "stat",
		Short: "Print stats for a license in the local relay server's pool",
		Run: func(cmd *cobra.Command, args []string) {
			license, err := manager.GetLicenseByID(cmd.Context(), licenseID)
			if err != nil {
				fmt.Printf("Error fetching stats for license ID %s: %v\n", licenseID, err)
				return
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
			if license.NodeID.Valid {
				nodeIDStr = strconv.FormatInt(license.NodeID.Int64, 10)
			} else {
				nodeIDStr = "-"
			}

			lastClaimedAtStr := formatTime(license.LastClaimedAt)
			lastReleasedAtStr := formatTime(license.LastReleasedAt)

			tableRows := []table.Row{
				{license.ID, claimsStr, nodeIDStr, lastClaimedAtStr, lastReleasedAtStr},
			}

			t := table.New(
				table.WithColumns(columns),
				table.WithRows(tableRows),
				table.WithFocused(true),
				table.WithHeight(5),
			)

			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(true)

			t.SetStyles(s)

			m := common.TableModel{Table: t}
			if _, err := tea.NewProgram(m).Run(); err != nil {
				fmt.Println("Error running program:", err)
			}
		},
	}

	cmd.Flags().StringVar(&licenseID, "id", "", "License ID to print stats for")
	cmd.MarkFlagRequired("id")

	return cmd
}
