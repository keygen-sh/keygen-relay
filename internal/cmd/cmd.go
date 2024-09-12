package cmd

import (
	"context"
	"errors"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os/exec"
)

var (
	dbPath    string
	verbosity int
)

func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	var rootCmd = &cobra.Command{
		Use:   "relay",
		Short: "Keygen Relay CLI",
	}

	rootCmd.PersistentFlags().StringVar(&dbPath, "database", "./relay.sqlite", "specify an alternate database path (default: ./relay.sqlite)")
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "verbosity level")

	conn, err := db.ConnectDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)

	rootCmd.AddCommand(addCmd(queries))

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}
	return 0
}
