package tui

import (
	"database/sql"
	"fmt"
	"strings"
)

type TableView struct {
	columns []string
	rows    [][]string
	cursor  int
	width   int
	height  int
}

func NewTableView(width, height int) TableView {
	return TableView{
		width:  width,
		height: height,
	}
}

func (t *TableView) Load(sqlRows *sql.Rows) error {
	cols, err := sqlRows.Columns()
	if err != nil {
		return fmt.Errorf("reading columns: %w", err)
	}
	t.columns = cols
	t.rows = nil
	t.cursor = 0

	for sqlRows.Next() {
		raw := make([]interface{}, len(cols))
		dest := make([]interface{}, len(cols))
		for i := range raw {
			dest[i] = &raw[i]
		}
		if err := sqlRows.Scan(dest...); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}
		row := make([]string, len(cols))
		for i, val := range raw {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		t.rows = append(t.rows, row)
	}
	return nil
}

func (t *TableView) MoveUp() {
	if t.cursor > 0 {
		t.cursor--
	}
}

func (t *TableView) MoveDown() {
	if t.cursor < len(t.rows)-1 {
		t.cursor++
	}
}

func (t TableView) Render() string {
	if len(t.columns) == 0 {
		return ErrorStyle.Render("\n  no data — select a table from the sidebar")
	}

	colW := t.colWidth()
	var b strings.Builder

	// header
	header := ""
	for _, col := range t.columns {
		header += TableHeaderStyle.Width(colW).Render(truncate(col, colW-2)) + " "
	}
	b.WriteString(header + "\n")

	// divider
	totalWidth := (colW+1)*len(t.columns)
	if totalWidth > t.width-28 {
		totalWidth = t.width - 28
	}
	b.WriteString(TableDividerStyle.Render(strings.Repeat("─", totalWidth)) + "\n")

	// visible rows
	visibleRows := t.height - 5
	if visibleRows < 1 {
		visibleRows = 1
	}
	start := 0
	if t.cursor >= visibleRows {
		start = t.cursor - visibleRows + 1
	}

	for i := start; i < len(t.rows) && i < start+visibleRows; i++ {
		row := t.rows[i]
		line := ""
		isSelected := i == t.cursor

		// row number indicator
		if isSelected {
			line += RowCursorStyle.Render("▶ ")
		} else {
			line += RowCursorStyle.Render("  ")
		}

		for _, cell := range row {
			if isSelected {
				line += TableCellActiveStyle.Width(colW).Render(truncate(cell, colW-2)) + " "
			} else {
				line += TableCellStyle.Width(colW).Render(truncate(cell, colW-2)) + " "
			}
		}
		b.WriteString(line + "\n")
	}

	// row counter at bottom
	if len(t.rows) > 0 {
		b.WriteString("\n" + RowCountStyle.Render(
			fmt.Sprintf("  row %d of %d", t.cursor+1, len(t.rows)),
		))
	}

	return b.String()
}

func (t TableView) colWidth() int {
	if len(t.columns) == 0 {
		return 20
	}
	available := t.width - 4
	w := available / len(t.columns)
	if w < 12 {
		return 12
	}
	if w > 24 {
		return 24
	}
	return w
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}