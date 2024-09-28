package output

import (
	"fmt"
	"github.com/fatih/color"
	"io"
)

func PrintError(w io.Writer, message string) {
	errorMessage := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s: %s\n", errorMessage("Error"), message)
}

func PrintSuccess(w io.Writer, message string) {
	successMessage := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s\n", successMessage(message))
}
