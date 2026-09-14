package ui

import (
	"fmt"
	"io"
	"strings"
)

type Column struct {
	Title string
	Width int
}

type Table struct {
	Columns []Column
	Rows    [][]string
	Theme   Theme
}

func NewTable(theme Theme, columns ...Column) *Table {
	return &Table{
		Columns: columns,
		Theme:   theme,
		Rows:    make([][]string, 0),
	}
}

func (t *Table) AddRow(cells ...string) {
	t.Rows = append(t.Rows, cells)
}

func (t *Table) Render(w io.Writer) error {
	var sb strings.Builder

	// Top Border
	sb.WriteString(TrueColor("┌", t.Theme.Border))
	for i, col := range t.Columns {
		sb.WriteString(TrueColor(strings.Repeat("─", col.Width+2), t.Theme.Border))
		if i < len(t.Columns)-1 {
			sb.WriteString(TrueColor("┬", t.Theme.Border))
		}
	}
	sb.WriteString(TrueColor("┐\n", t.Theme.Border))

	// Header Row
	sb.WriteString(TrueColor("│", t.Theme.Border))
	for _, col := range t.Columns {
		title := fmt.Sprintf(" %-*s ", col.Width, col.Title)
		sb.WriteString(TrueColor(title, t.Theme.Primary))
		sb.WriteString(TrueColor("│", t.Theme.Border))
	}
	sb.WriteString("\n")

	// Header Divider
	sb.WriteString(TrueColor("├", t.Theme.Border))
	for i, col := range t.Columns {
		sb.WriteString(TrueColor(strings.Repeat("─", col.Width+2), t.Theme.Border))
		if i < len(t.Columns)-1 {
			sb.WriteString(TrueColor("┼", t.Theme.Border))
		}
	}
	sb.WriteString(TrueColor("┤\n", t.Theme.Border))

	// Data Rows
	for _, row := range t.Rows {
		sb.WriteString(TrueColor("│", t.Theme.Border))
		for i, cell := range row {
			width := t.Columns[i].Width
			if len(cell) > width {
				cell = cell[:width-3] + "..."
			}
			padded := fmt.Sprintf(" %-*s ", width, cell)
			sb.WriteString(TrueColor(padded, t.Theme.Accent))
			sb.WriteString(TrueColor("│", t.Theme.Border))
		}
		sb.WriteString("\n")
	}

	// Bottom Border
	sb.WriteString(TrueColor("└", t.Theme.Border))
	for i, col := range t.Columns {
		sb.WriteString(TrueColor(strings.Repeat("─", col.Width+2), t.Theme.Border))
		if i < len(t.Columns)-1 {
			sb.WriteString(TrueColor("┴", t.Theme.Border))
		}
	}
	sb.WriteString(TrueColor("┘\n", t.Theme.Border))

	_, err := fmt.Fprint(w, sb.String())
	return err
}
