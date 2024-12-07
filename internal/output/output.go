package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

func PrintError(w io.Writer, message string) {
	errorMessage := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s: %s\n", errorMessage("error"), message)
}

func PrintSuccess(w io.Writer, message string) {
	successMessage := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s\n", successMessage(message))
}

func Print(w io.Writer, message string) {
	fmt.Fprintf(w, "%s\n", message)
}
