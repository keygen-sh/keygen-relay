package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

func PrintError(w io.Writer, message string, v ...any) {
	errorMessage := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s: %s\n", errorMessage("error"), fmt.Sprintf(message, v...))
}

func PrintSuccess(w io.Writer, message string, v ...any) {
	successMessage := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Fprintf(w, "%s\n", successMessage(fmt.Sprintf(message, v...)))
}

func Print(w io.Writer, message string, v ...any) {
	fmt.Fprintf(w, "%s\n", fmt.Sprintf(message, v...))
}
