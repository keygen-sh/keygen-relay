package ui

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/bubbles/table"
)

type SimpleTableRenderer struct {
	Output io.Writer
}

func NewSimpleTableRenderer(output io.Writer) TableRenderer {
	return &SimpleTableRenderer{
		Output: output,
	}
}

func (r *SimpleTableRenderer) Render(rows []table.Row, columns []table.Column) error {
	w := tabwriter.NewWriter(r.Output, 0, 0, 2, ' ', 0)

	// determine the maximum width of each column
	maxWidths := make([]int, len(columns))

	// start with the length of the headers
	for i, col := range columns {
		maxWidths[i] = len(col.Title)
	}

	for _, row := range rows {
		for i, cell := range row {
			cellLength := len(fmt.Sprintf("%v", cell))
			if cellLength > maxWidths[i] {
				maxWidths[i] = cellLength
			}
		}
	}

	// render the headers
	for i, col := range columns {
		suffix := strings.Repeat(" ", maxWidths[i]-len(col.Title)) + " |"
		prefix := " "

		if i == 0 {
			prefix = "| "
		}

		fmt.Fprintf(w, "%s%v%s", prefix, col.Title, suffix)
	}
	fmt.Fprintln(w)

	// render the lines
	for i, width := range maxWidths {
		suffix := "+"
		prefix := ""

		if i == 0 {
			prefix = "+"
		}

		fmt.Fprint(w, prefix+strings.Repeat("-", width+2)+suffix)
	}
	fmt.Fprintln(w)

	// render the rows
	for _, row := range rows {
		for i, cell := range row {
			suffix := strings.Repeat(" ", maxWidths[i]-len(cell)) + " |"
			prefix := " "

			if i == 0 {
				prefix = "| "
			}

			fmt.Fprintf(w, "%s%v%s", prefix, cell, suffix)
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
}
