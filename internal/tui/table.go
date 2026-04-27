package tui

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/HalxDocs/lazydb/internal/db"
)

type SortDir int

const (
	SortNone SortDir = iota
	SortAsc
	SortDesc
)

type cellKind int

const (
	kindString cellKind = iota
	kindNumeric
	kindBool
	kindNull
)

type Row struct {
	Cells []string
	Kinds []cellKind
}

type TableView struct {
	columns     []string
	rows        []Row
	allRows     []Row // unfiltered + unsorted snapshot
	rowCursor   int
	colCursor   int
	colOffset   int // horizontal scroll: index of leftmost rendered column
	width       int
	height      int
	sortCol     int
	sortDir     SortDir
	filter      string
	loadedAt    time.Time
	queryMS     int64
	schema      []db.Column
	showSchema  bool
	tableName   string
	rowsLimit   int
	totalCount  int // total in DB (best-effort), -1 = unknown
}

func NewTableView(width, height int) TableView {
	return TableView{
		width:      width,
		height:     height,
		sortCol:    -1,
		rowsLimit:  500,
		totalCount: -1,
	}
}

func (t *TableView) SetSize(w, h int) {
	t.width = w
	t.height = h
}

func (t *TableView) SetTable(name string, schema []db.Column, total int) {
	t.tableName = name
	t.schema = schema
	t.totalCount = total
	t.showSchema = false
	t.filter = ""
	t.sortCol = -1
	t.sortDir = SortNone
	t.colCursor = 0
	t.colOffset = 0
}

func (t *TableView) Load(sqlRows *sql.Rows, took time.Duration) error {
	cols, err := sqlRows.Columns()
	if err != nil {
		return fmt.Errorf("reading columns: %w", err)
	}
	colTypes, _ := sqlRows.ColumnTypes()
	_ = colTypes

	t.columns = cols
	t.allRows = nil
	t.rowCursor = 0
	if t.colCursor >= len(cols) {
		t.colCursor = 0
	}

	for sqlRows.Next() {
		raw := make([]any, len(cols))
		dest := make([]any, len(cols))
		for i := range raw {
			dest[i] = &raw[i]
		}
		if err := sqlRows.Scan(dest...); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}
		row := Row{
			Cells: make([]string, len(cols)),
			Kinds: make([]cellKind, len(cols)),
		}
		for i, val := range raw {
			row.Cells[i], row.Kinds[i] = formatValue(val)
		}
		t.allRows = append(t.allRows, row)
	}
	t.queryMS = took.Milliseconds()
	t.loadedAt = time.Now()
	t.applyFilterSort()
	return nil
}

func formatValue(val any) (string, cellKind) {
	if val == nil {
		return "NULL", kindNull
	}
	switch v := val.(type) {
	case bool:
		if v {
			return "true", kindBool
		}
		return "false", kindBool
	case []byte:
		s := string(v)
		if isNumeric(s) {
			return s, kindNumeric
		}
		return s, kindString
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v), kindNumeric
	case time.Time:
		return v.Format("2006-01-02 15:04:05"), kindString
	case string:
		if isNumeric(v) {
			return v, kindNumeric
		}
		return v, kindString
	default:
		s := fmt.Sprintf("%v", v)
		if isNumeric(s) {
			return s, kindNumeric
		}
		return s, kindString
	}
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil
}

// ---- navigation ----

func (t *TableView) MoveUp() {
	if t.rowCursor > 0 {
		t.rowCursor--
	}
}

func (t *TableView) MoveDown() {
	if t.rowCursor < len(t.rows)-1 {
		t.rowCursor++
	}
}

func (t *TableView) MoveLeft() {
	if t.colCursor > 0 {
		t.colCursor--
		if t.colCursor < t.colOffset {
			t.colOffset = t.colCursor
		}
	}
}

func (t *TableView) MoveRight() {
	if t.colCursor < len(t.columns)-1 {
		t.colCursor++
	}
}

func (t *TableView) PageDown() {
	step := t.visibleRowCount()
	t.rowCursor += step
	if t.rowCursor >= len(t.rows) {
		t.rowCursor = len(t.rows) - 1
	}
	if t.rowCursor < 0 {
		t.rowCursor = 0
	}
}

func (t *TableView) PageUp() {
	step := t.visibleRowCount()
	t.rowCursor -= step
	if t.rowCursor < 0 {
		t.rowCursor = 0
	}
}

func (t *TableView) Top()    { t.rowCursor = 0 }
func (t *TableView) Bottom() { t.rowCursor = max(0, len(t.rows)-1) }

func (t *TableView) RowCursor() int     { return t.rowCursor }
func (t *TableView) ColCursor() int     { return t.colCursor }
func (t *TableView) RowCount() int      { return len(t.rows) }
func (t *TableView) Columns() []string  { return t.columns }
func (t *TableView) Schema() []db.Column { return t.schema }
func (t *TableView) Filter() string     { return t.filter }
func (t *TableView) ShowingSchema() bool { return t.showSchema }
func (t *TableView) TableName() string  { return t.tableName }
func (t *TableView) QueryMS() int64     { return t.queryMS }
func (t *TableView) TotalCount() int    { return t.totalCount }

func (t *TableView) ToggleSchema() { t.showSchema = !t.showSchema }

func (t *TableView) CurrentRow() (Row, bool) {
	if t.rowCursor < 0 || t.rowCursor >= len(t.rows) {
		return Row{}, false
	}
	return t.rows[t.rowCursor], true
}

func (t *TableView) CurrentCell() (string, bool) {
	r, ok := t.CurrentRow()
	if !ok {
		return "", false
	}
	if t.colCursor < 0 || t.colCursor >= len(r.Cells) {
		return "", false
	}
	return r.Cells[t.colCursor], true
}

// SetFilter sets a substring filter applied across all cells.
func (t *TableView) SetFilter(f string) {
	t.filter = f
	t.applyFilterSort()
}

// CycleSort cycles asc → desc → none for the current column.
func (t *TableView) CycleSort() {
	if t.colCursor < 0 || t.colCursor >= len(t.columns) {
		return
	}
	if t.sortCol != t.colCursor {
		t.sortCol = t.colCursor
		t.sortDir = SortAsc
	} else {
		switch t.sortDir {
		case SortNone:
			t.sortDir = SortAsc
		case SortAsc:
			t.sortDir = SortDesc
		case SortDesc:
			t.sortDir = SortNone
			t.sortCol = -1
		}
	}
	t.applyFilterSort()
}

func (t *TableView) applyFilterSort() {
	// filter
	t.rows = t.rows[:0]
	if t.filter == "" {
		t.rows = append(t.rows, t.allRows...)
	} else {
		needle := strings.ToLower(t.filter)
		for _, r := range t.allRows {
			for _, c := range r.Cells {
				if strings.Contains(strings.ToLower(c), needle) {
					t.rows = append(t.rows, r)
					break
				}
			}
		}
	}
	// sort
	if t.sortCol >= 0 && t.sortCol < len(t.columns) && t.sortDir != SortNone {
		col := t.sortCol
		dir := t.sortDir
		sort.SliceStable(t.rows, func(i, j int) bool {
			a := ""
			b := ""
			if col < len(t.rows[i].Cells) {
				a = t.rows[i].Cells[col]
			}
			if col < len(t.rows[j].Cells) {
				b = t.rows[j].Cells[col]
			}
			less := compareCells(a, b)
			if dir == SortDesc {
				return !less
			}
			return less
		})
	}
	if t.rowCursor >= len(t.rows) {
		t.rowCursor = max(0, len(t.rows)-1)
	}
}

func compareCells(a, b string) bool {
	an, aerr := strconv.ParseFloat(strings.TrimSpace(a), 64)
	bn, berr := strconv.ParseFloat(strings.TrimSpace(b), 64)
	if aerr == nil && berr == nil {
		return an < bn
	}
	return strings.ToLower(a) < strings.ToLower(b)
}

// ExportCSV writes the currently visible rows to a CSV file under the user's
// home dir and returns the absolute path.
func (t *TableView) ExportCSV() (string, error) {
	if len(t.columns) == 0 {
		return "", fmt.Errorf("no data to export")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".lazydb", "exports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := t.tableName
	if name == "" {
		name = "query"
	}
	stamp := time.Now().Format("20060102-150405")
	path := filepath.Join(dir, fmt.Sprintf("%s-%s.csv", sanitizeFilename(name), stamp))
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write(t.columns); err != nil {
		return "", err
	}
	for _, r := range t.rows {
		if err := w.Write(r.Cells); err != nil {
			return "", err
		}
	}
	w.Flush()
	return path, w.Error()
}

func sanitizeFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// ---- rendering ----

func (t *TableView) visibleRowCount() int {
	r := t.height - 4 // header + divider + footer
	if r < 1 {
		return 1
	}
	return r
}

// keep colOffset valid given current cursor and rendered width
func (t *TableView) ensureColVisible(colW int) {
	if len(t.columns) == 0 {
		return
	}
	maxCols := (t.width - 4) / (colW + 1)
	if maxCols < 1 {
		maxCols = 1
	}
	if t.colCursor < t.colOffset {
		t.colOffset = t.colCursor
	}
	if t.colCursor >= t.colOffset+maxCols {
		t.colOffset = t.colCursor - maxCols + 1
	}
	if t.colOffset < 0 {
		t.colOffset = 0
	}
	if t.colOffset > len(t.columns)-1 {
		t.colOffset = max(0, len(t.columns)-1)
	}
}

func (t TableView) Render() string {
	if t.showSchema {
		return t.renderSchema()
	}
	if len(t.columns) == 0 {
		hint := "  no data — pick a table from the sidebar  ·  press " +
			KbdStyle.Render("?") + " for help"
		return "\n" + hint
	}

	colW := t.colWidth()
	(&t).ensureColVisible(colW)

	maxCols := (t.width - 4) / (colW + 1)
	if maxCols < 1 {
		maxCols = 1
	}
	end := t.colOffset + maxCols
	if end > len(t.columns) {
		end = len(t.columns)
	}

	var b strings.Builder

	// header
	b.WriteString("  ")
	for i := t.colOffset; i < end; i++ {
		col := t.columns[i]
		label := truncate(col, colW-3)
		// sort indicator
		if i == t.sortCol && t.sortDir != SortNone {
			ind := SortAscIndicator
			if t.sortDir == SortDesc {
				ind = SortDescIndicator
			}
			label = label + " " + ind
		}
		if i == t.colCursor {
			b.WriteString(TableHeaderActiveStyle.Width(colW).Render(" " + label) + " ")
		} else {
			b.WriteString(TableHeaderStyle.Width(colW).Render(" " + label) + " ")
		}
	}
	if end < len(t.columns) {
		b.WriteString(BreadcrumbSepStyle.Render("→"))
	}
	b.WriteString("\n")

	// divider
	totalW := (colW+1)*(end-t.colOffset) + 2
	if totalW > t.width {
		totalW = t.width
	}
	b.WriteString(TableDividerStyle.Render(strings.Repeat("─", totalW)) + "\n")

	// rows
	visible := t.visibleRowCount()
	start := 0
	if t.rowCursor >= visible {
		start = t.rowCursor - visible + 1
	}
	stop := start + visible
	if stop > len(t.rows) {
		stop = len(t.rows)
	}

	for i := start; i < stop; i++ {
		row := t.rows[i]
		isSelectedRow := i == t.rowCursor

		if isSelectedRow {
			b.WriteString(RowCursorStyle.Render("▶ "))
		} else {
			b.WriteString("  ")
		}

		for c := t.colOffset; c < end; c++ {
			cell := ""
			kind := kindString
			if c < len(row.Cells) {
				cell = row.Cells[c]
				kind = row.Kinds[c]
			}
			text := truncate(cell, colW-2)
			isFocused := isSelectedRow && c == t.colCursor

			switch {
			case isFocused:
				b.WriteString(TableCellFocusedStyle.Width(colW).Render(" " + text) + " ")
			case isSelectedRow:
				st := TableCellActiveStyle.Width(colW)
				if kind == kindNumeric {
					st = st.Foreground(colorNumeric)
				} else if kind == kindNull {
					st = st.Foreground(colorNull).Italic(true)
				} else if kind == kindBool {
					st = st.Foreground(colorBool)
				}
				b.WriteString(st.Render(" " + text) + " ")
			default:
				st := TableCellStyle.Width(colW)
				switch kind {
				case kindNumeric:
					st = TableCellNumericStyle.Width(colW)
				case kindNull:
					st = TableCellNullStyle.Width(colW)
				case kindBool:
					st = TableCellBoolStyle.Width(colW)
				}
				b.WriteString(st.Render(" " + text) + " ")
			}
		}
		b.WriteString("\n")
	}

	// pad blank rows so the layout stays stable
	for i := stop - start; i < visible; i++ {
		b.WriteString("\n")
	}

	// footer line
	footer := t.renderFooter()
	b.WriteString(footer)

	return b.String()
}

func (t TableView) renderFooter() string {
	if len(t.rows) == 0 {
		if t.filter != "" {
			return "  " + ToastErrorStyle.Render(" no rows match filter ") + " "
		}
		return ""
	}
	parts := []string{}
	parts = append(parts, fmt.Sprintf("row %d/%d", t.rowCursor+1, len(t.rows)))
	if t.totalCount >= 0 && t.totalCount > len(t.allRows) {
		parts = append(parts, fmt.Sprintf("of %d total", t.totalCount))
	}
	if t.filter != "" {
		parts = append(parts, fmt.Sprintf("filter: %q", t.filter))
	}
	if t.sortCol >= 0 && t.sortCol < len(t.columns) && t.sortDir != SortNone {
		dir := SortAscIndicator
		if t.sortDir == SortDesc {
			dir = SortDescIndicator
		}
		parts = append(parts, fmt.Sprintf("sort: %s %s", t.columns[t.sortCol], dir))
	}
	if t.queryMS > 0 {
		parts = append(parts, fmt.Sprintf("%dms", t.queryMS))
	}
	return RowCountStyle.Render(strings.Join(parts, "  ·  "))
}

func (t TableView) renderSchema() string {
	if len(t.schema) == 0 {
		return "\n  schema unavailable"
	}
	var b strings.Builder
	title := fmt.Sprintf("  schema: %s", t.tableName)
	b.WriteString(BreadcrumbActiveStyle.Render(title) + "\n\n")

	colNameW := 24
	colTypeW := 20
	header := fmt.Sprintf("  %-*s %-*s %s",
		colNameW, "column",
		colTypeW, "type",
		"nullable")
	b.WriteString(BreadcrumbStyle.Render(header) + "\n")
	b.WriteString(TableDividerStyle.Render(strings.Repeat("─", colNameW+colTypeW+18)) + "\n")

	for _, c := range t.schema {
		nullable := "NOT NULL"
		nullStyle := SchemaPKStyle
		if c.Nullable {
			nullable = "nullable"
			nullStyle = SchemaNullableStyle
		}
		line := fmt.Sprintf("  %-*s ", colNameW, truncate(c.Name, colNameW))
		line += SchemaTypeStyle.Render(fmt.Sprintf("%-*s ", colTypeW, truncate(c.Type, colTypeW)))
		line += nullStyle.Render(nullable)
		b.WriteString(line + "\n")
	}
	b.WriteString("\n  " + ModalDimStyle.Render("press ") + KbdStyle.Render("s") +
		ModalDimStyle.Render(" to return to data view"))
	return b.String()
}

func (t TableView) colWidth() int {
	if len(t.columns) == 0 {
		return 20
	}
	available := t.width - 4
	target := 22
	w := available / max(1, min(len(t.columns), available/12))
	_ = target
	if w < 12 {
		w = 12
	}
	if w > 36 {
		w = 36
	}
	return w
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	// strip newlines for table cells
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
