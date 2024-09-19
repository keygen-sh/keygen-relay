package ui

import "github.com/charmbracelet/bubbles/table"

type TableRenderer interface {
	Render(rows []table.Row, columns []table.Column) error
}
