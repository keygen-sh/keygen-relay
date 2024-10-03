package main

import (
	"os"

	"github.com/keygen-sh/keygen-relay/cli"
)

func main() {
	os.Exit(cli.Run())
}
