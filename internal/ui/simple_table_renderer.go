package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"io"
	"strings"
	"text/tabwriter"
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

	for _, col := range columns {
		fmt.Fprintf(w, "%s\t", col.Title)
	}
	fmt.Fprintln(w)

	for range columns {
		fmt.Fprintf(w, "%s\t", strings.Repeat("-", 15))
	}
	fmt.Fprintln(w)

	for _, row := range rows {
		for _, cell := range row {
			fmt.Fprintf(w, "%v\t", cell)
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
}
