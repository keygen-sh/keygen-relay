package cmd

import (
	"github.com/spf13/cobra"
)

func HelpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "help [command]",
		Short:        "help for a command",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd, _, err := c.Root().Find(args)
			if cmd == nil || err != nil {
				return err
			}

			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
