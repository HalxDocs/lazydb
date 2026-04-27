package tui

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/HalxDocs/lazydb/internal/db"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeQuery        // single-line query
	ModeFilter       // single-line filter
	ModeEditor       // multi-line SQL editor
	ModeHelp
	ModeDetail
)

type Model struct {
	db        db.DB
	driver    string
	sidebar   Sidebar
	tableView TableView
	queryBar  QueryBar
	editor    SQLEditor
	mode      Mode
	width     int
	height    int
	err       error
	status    string
	ready     bool
	loading   bool
	spinner   spinner.Model
	toast     toastMsg
	startedAt time.Time
}

type toastMsg struct {
	text   string
	until  time.Time
	isErr  bool
}

// ---- bubbletea messages ----

type tablesLoadedMsg struct{ tables []TableMeta }
type rowsLoadedMsg struct {
	rows  *sql.Rows
	took  time.Duration
	table string
	count int
	cols  []db.Column
}
type queryRanMsg struct {
	rows  *sql.Rows
	took  time.Duration
	query string
}
type errMsg struct{ err error }
type clearToastMsg struct{}
type tickMsg time.Time

func (e errMsg) Error() string { return e.err.Error() }

// ---- constructors ----

func NewModel(database db.DB) Model {
	return NewModelWithDriver(database, "")
}

func NewModelWithDriver(database db.DB, driver string) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorAccent)

	m := Model{
		db:        database,
		driver:    driver,
		queryBar:  NewQueryBar(),
		editor:    NewSQLEditor(),
		spinner:   sp,
		startedAt: time.Now(),
	}
	m.sidebar.SetDriver(driver)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadTables(), m.spinner.Tick)
}

// ---- update ----

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			sidebarW, bodyH := m.layoutSizes()
			mainW := m.width - sidebarW
			m.sidebar = NewSidebar(sidebarW, bodyH)
			m.sidebar.SetDriver(m.driver)
			m.tableView = NewTableView(mainW, bodyH)
		}
		m.recomputeLayout()
		m.ready = true
		if len(m.sidebar.Tables()) == 0 {
			return m, m.loadTables()
		}
		return m, nil

	case tablesLoadedMsg:
		m.sidebar.SetTables(msg.tables)
		m.loading = false
		m.status = fmt.Sprintf("%d tables loaded in %s",
			len(msg.tables), time.Since(m.startedAt).Round(time.Millisecond))
		if len(msg.tables) > 0 {
			return m, m.loadRows(msg.tables[0].Name)
		}
		return m, nil

	case rowsLoadedMsg:
		if msg.rows != nil {
			m.tableView.SetTable(msg.table, msg.cols, msg.count)
			if err := m.tableView.Load(msg.rows, msg.took); err != nil {
				m.err = err
			}
			msg.rows.Close()
		}
		m.loading = false
		return m, nil

	case queryRanMsg:
		if msg.rows != nil {
			m.tableView.SetTable("query", nil, -1)
			if err := m.tableView.Load(msg.rows, msg.took); err != nil {
				m.err = err
			}
			msg.rows.Close()
			m.queryBar.PushHistory(msg.query)
			return m, m.toastCmd(fmt.Sprintf("query · %d rows · %dms",
				m.tableView.RowCount(), msg.took.Milliseconds()), false)
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, m.toastCmd(msg.err.Error(), true)

	case clearToastMsg:
		if time.Now().After(m.toast.until) {
			m.toast = toastMsg{}
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// passthrough for inputs
	if m.mode == ModeQuery || m.mode == ModeFilter {
		var cmd tea.Cmd
		m.queryBar.input, cmd = m.queryBar.input.Update(msg)
		if m.mode == ModeFilter {
			m.tableView.SetFilter(m.queryBar.Value())
		}
		return m, cmd
	}
	if m.mode == ModeEditor {
		var cmd tea.Cmd
		m.editor.area, cmd = m.editor.area.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// modal modes first
	switch m.mode {
	case ModeQuery:
		return m.handleQueryKey(msg)
	case ModeFilter:
		return m.handleFilterKey(msg)
	case ModeEditor:
		return m.handleEditorKey(msg)
	case ModeHelp:
		return m.handleHelpKey(msg)
	case ModeDetail:
		return m.handleDetailKey(msg)
	}
	return m.handleNormalKey(msg)
}

func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		switch msg.String() {
		case "esc", "enter":
			m.err = nil
			return m, nil
		}
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.mode = ModeHelp
		return m, nil

	case "up", "k":
		m.tableView.MoveUp()

	case "down", "j":
		m.tableView.MoveDown()

	case "left", "h":
		m.tableView.MoveLeft()

	case "right", "l":
		m.tableView.MoveRight()

	case "g":
		m.tableView.Top()
	case "G":
		m.tableView.Bottom()
	case "ctrl+d":
		m.tableView.PageDown()
	case "ctrl+u":
		m.tableView.PageUp()

	case "tab", "]":
		m.sidebar.MoveDown()
		if t := m.sidebar.SelectedTable(); t != "" {
			m.loading = true
			return m, m.loadRows(t)
		}

	case "shift+tab", "[":
		m.sidebar.MoveUp()
		if t := m.sidebar.SelectedTable(); t != "" {
			m.loading = true
			return m, m.loadRows(t)
		}

	case "K":
		m.sidebar.MoveUp()
		if t := m.sidebar.SelectedTable(); t != "" {
			m.loading = true
			return m, m.loadRows(t)
		}

	case "J":
		m.sidebar.MoveDown()
		if t := m.sidebar.SelectedTable(); t != "" {
			m.loading = true
			return m, m.loadRows(t)
		}

	case "enter":
		if m.tableView.RowCount() > 0 {
			m.mode = ModeDetail
		}

	case "s":
		m.tableView.ToggleSchema()

	case "o":
		m.tableView.CycleSort()

	case "f":
		m.mode = ModeFilter
		filteringSidebar = false
		m.queryBar.ShowFilter(m.tableView.Filter())
		m.recomputeLayout()

	case "r":
		if t := m.sidebar.SelectedTable(); t != "" {
			m.loading = true
			return m, m.loadRows(t)
		}

	case "/":
		m.mode = ModeQuery
		m.queryBar.ShowSQL()
		m.recomputeLayout()

	case ":":
		m.mode = ModeEditor
		m.editor.Resize(min(100, m.width-10))
		m.editor.Show()

	case "t":
		// search tables — reuse query bar in filter mode but for sidebar.
		// simplest: open filter bar with sidebar prefix marker.
		m.mode = ModeFilter
		m.queryBar.ShowFilter(m.sidebar.Search())
		// route to sidebar instead — set a flag via mode? Use prefix trick:
		// we'll piggy-back: use queryBar in filter mode but check a marker.
		// To keep things tidy, we set a sentinel state on sidebar by using
		// the same filter bar; final value is committed on Enter.
		m.queryBar.input.Placeholder = "search tables (esc to cancel)"
		m.queryBar.input.SetValue(m.sidebar.Search())
		m.sidebar.SetSearch(m.queryBar.Value())
		// We treat 't' filter as sidebar search — track with an internal flag:
		filteringSidebar = true
		m.recomputeLayout()
		return m, nil

	case "e":
		path, err := m.tableView.ExportCSV()
		if err != nil {
			return m, m.toastCmd("export failed: "+err.Error(), true)
		}
		return m, m.toastCmd("exported → "+path, false)

	case "y":
		if v, ok := m.tableView.CurrentCell(); ok {
			if err := clipboard.WriteAll(v); err != nil {
				return m, m.toastCmd("clipboard error: "+err.Error(), true)
			}
			return m, m.toastCmd("yanked cell ("+strconvLen(v)+" chars)", false)
		}

	case "Y":
		if r, ok := m.tableView.CurrentRow(); ok {
			tsv := strings.Join(r.Cells, "\t")
			if err := clipboard.WriteAll(tsv); err != nil {
				return m, m.toastCmd("clipboard error: "+err.Error(), true)
			}
			return m, m.toastCmd("yanked row as TSV", false)
		}

	case "esc":
		m.err = nil
	}

	return m, nil
}

// strconvLen avoids pulling in strconv just for itoa.
func strconvLen(s string) string {
	return fmt.Sprintf("%d", len(s))
}

// filteringSidebar tracks if the current filter bar is searching the sidebar
// rather than table rows. Module-scoped because the filter bar is a single
// shared widget.
var filteringSidebar bool

func (m Model) handleQueryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.queryBar.Hide()
		m.recomputeLayout()
		return m, nil

	case "up":
		m.queryBar.HistoryUp()
		return m, nil

	case "down":
		m.queryBar.HistoryDown()
		return m, nil

	case "enter":
		query := strings.TrimSpace(m.queryBar.Value())
		m.mode = ModeNormal
		m.queryBar.Hide()
		m.recomputeLayout()
		if query != "" {
			m.loading = true
			return m, m.runQuery(query)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.queryBar.input, cmd = m.queryBar.input.Update(msg)
	return m, cmd
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		if filteringSidebar {
			m.sidebar.SetSearch("")
			filteringSidebar = false
		} else {
			m.tableView.SetFilter("")
		}
		m.queryBar.Hide()
		m.recomputeLayout()
		return m, nil
	case "enter":
		m.mode = ModeNormal
		if filteringSidebar {
			m.sidebar.SetSearch(m.queryBar.Value())
			filteringSidebar = false
			if t := m.sidebar.SelectedTable(); t != "" {
				m.loading = true
				m.queryBar.Hide()
				m.recomputeLayout()
				return m, m.loadRows(t)
			}
		} else {
			m.tableView.SetFilter(m.queryBar.Value())
		}
		m.queryBar.Hide()
		m.recomputeLayout()
		return m, nil
	}

	var cmd tea.Cmd
	m.queryBar.input, cmd = m.queryBar.input.Update(msg)
	if filteringSidebar {
		m.sidebar.SetSearch(m.queryBar.Value())
	} else {
		m.tableView.SetFilter(m.queryBar.Value())
	}
	return m, cmd
}

func (m Model) handleEditorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.editor.Hide()
		return m, nil
	case "ctrl+enter", "alt+enter":
		query := strings.TrimSpace(m.editor.Value())
		m.mode = ModeNormal
		m.editor.Hide()
		if query != "" {
			m.loading = true
			return m, m.runQuery(query)
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.editor.area, cmd = m.editor.area.Update(msg)
	return m, cmd
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?", "q":
		m.mode = ModeNormal
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q":
		m.mode = ModeNormal
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "y":
		if v, ok := m.tableView.CurrentCell(); ok {
			if err := clipboard.WriteAll(v); err != nil {
				return m, m.toastCmd("clipboard error: "+err.Error(), true)
			}
			return m, m.toastCmd("yanked cell", false)
		}
	case "up", "k":
		m.tableView.MoveUp()
	case "down", "j":
		m.tableView.MoveDown()
	}
	return m, nil
}

// ---- view ----

func (m Model) View() string {
	if !m.ready {
		return "\n  " + m.spinner.View() + " connecting to database..."
	}

	body := m.renderBody()

	switch m.mode {
	case ModeHelp:
		return renderHelp(m.width, m.height)
	case ModeDetail:
		if r, ok := m.tableView.CurrentRow(); ok {
			return renderDetail(m.tableView.Columns(), r, m.width, m.height)
		}
	case ModeEditor:
		return m.editor.Render(m.width, m.height)
	}

	queryBar := ""
	if m.queryBar.IsVisible() {
		queryBar = m.queryBar.Render()
	}

	statusBar := m.renderStatusBar()
	breadcrumb := m.renderBreadcrumb()

	parts := []string{breadcrumb, body}
	if queryBar != "" {
		parts = append(parts, queryBar)
	}
	parts = append(parts, statusBar)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) renderBody() string {
	if m.err != nil {
		errBox := ToastErrorStyle.Render(" ERROR ") + " " +
			lipgloss.NewStyle().Foreground(colorError).Render(m.err.Error())
		hint := ModalDimStyle.Render("press esc or enter to dismiss  ·  ? for help")
		return "\n  " + errBox + "\n\n  " + hint
	}

	sidebar := m.sidebar.Render()
	mainStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
	main := mainStyle.Render(m.tableView.Render())
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
}

func (m Model) renderBreadcrumb() string {
	driver := m.driver
	if driver == "" {
		driver = "db"
	}
	parts := []string{
		BreadcrumbActiveStyle.Render(" lazydb "),
		BreadcrumbSepStyle.Render(" / "),
		BreadcrumbStyle.Render(driver),
	}
	if t := m.sidebar.SelectedTable(); t != "" {
		parts = append(parts,
			BreadcrumbSepStyle.Render(" / "),
			BreadcrumbActiveStyle.Render(t),
		)
		if m.tableView.ShowingSchema() {
			parts = append(parts,
				BreadcrumbSepStyle.Render(" · "),
				BreadcrumbStyle.Render("schema"),
			)
		}
	}
	left := strings.Join(parts, "")
	right := ""
	if m.loading {
		right = BreadcrumbStyle.Render(m.spinner.View() + " loading ")
	}
	gap := lipgloss.NewStyle().
		Background(colorSurface).
		Width(max(0, m.width-lipgloss.Width(left)-lipgloss.Width(right))).
		Render("")
	return ToolbarStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right),
	)
}

func (m Model) renderStatusBar() string {
	mode := StatusBarModeStyle.Render(" NORMAL ")
	switch m.mode {
	case ModeQuery:
		mode = StatusBarModeQueryStyle.Render(" QUERY ")
	case ModeFilter:
		mode = StatusBarModeFilterStyle.Render(" FILTER ")
	case ModeEditor:
		mode = StatusBarModeQueryStyle.Render(" EDITOR ")
	case ModeHelp:
		mode = StatusBarModeStyle.Render(" HELP ")
	case ModeDetail:
		mode = StatusBarModeStyle.Render(" DETAIL ")
	}

	left := mode
	mid := ""
	if m.toast.text != "" && time.Now().Before(m.toast.until) {
		if m.toast.isErr {
			mid = "  " + ToastErrorStyle.Render(" "+m.toast.text+" ")
		} else {
			mid = "  " + ToastStyle.Render(" "+m.toast.text+" ")
		}
	} else {
		stat := ""
		if t, ok := m.sidebar.SelectedMeta(); ok {
			stat = fmt.Sprintf(" %s", t.Name)
			if t.Count > 0 {
				stat += fmt.Sprintf(" · %s rows", humanCount(t.Count))
			}
			if ms := m.tableView.QueryMS(); ms > 0 {
				stat += fmt.Sprintf(" · %dms", ms)
			}
		}
		mid = StatusBarTextStyle.Render(stat)
	}

	right := StatusBarHintStyle.Render("? help  /:query  f filter  o sort  e export  q quit ")

	leftW := lipgloss.Width(left)
	midW := lipgloss.Width(mid)
	rightW := lipgloss.Width(right)
	gap := m.width - leftW - midW - rightW
	if gap < 0 {
		gap = 0
	}
	spacer := lipgloss.NewStyle().
		Background(lipgloss.Color("#0A2040")).
		Width(gap).
		Render("")
	return StatusBarStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, mid, spacer, right),
	)
}

// ---- commands ----

func (m Model) loadTables() tea.Cmd {
	return func() tea.Msg {
		names, err := m.db.Tables()
		if err != nil {
			return errMsg{err}
		}
		// parallelize CountRows — much faster on databases with many tables
		tables := make([]TableMeta, len(names))
		var wg sync.WaitGroup
		sem := make(chan struct{}, 8)
		for i, name := range names {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int, name string) {
				defer wg.Done()
				defer func() { <-sem }()
				count, _ := m.db.CountRows(name)
				tables[i] = TableMeta{Name: name, Count: count}
			}(i, name)
		}
		wg.Wait()
		return tablesLoadedMsg{tables}
	}
}

func (m Model) loadRows(table string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		rows, err := m.db.Rows(table, 500)
		if err != nil {
			return errMsg{err}
		}
		cols, _ := m.db.Columns(table)
		count, _ := m.db.CountRows(table)
		return rowsLoadedMsg{
			rows:  rows,
			took:  time.Since(start),
			table: table,
			count: count,
			cols:  cols,
		}
	}
}

func (m Model) runQuery(query string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		rows, err := m.db.Query(query)
		if err != nil {
			return errMsg{err}
		}
		return queryRanMsg{
			rows:  rows,
			took:  time.Since(start),
			query: query,
		}
	}
}

// layoutSizes returns the sidebar width and the available body height,
// accounting for the breadcrumb, status bar, and (when visible) query bar.
func (m Model) layoutSizes() (sidebarW, bodyH int) {
	sidebarW = 30
	if m.width < 100 {
		sidebarW = 24
	}
	if m.width < 70 {
		sidebarW = 20
	}
	chrome := 2 // breadcrumb + status
	if m.queryBar.IsVisible() {
		chrome++
	}
	bodyH = m.height - chrome
	if bodyH < 5 {
		bodyH = 5
	}
	return sidebarW, bodyH
}

func (m *Model) recomputeLayout() {
	sidebarW, bodyH := m.layoutSizes()
	mainW := m.width - sidebarW
	m.sidebar.SetSize(sidebarW, bodyH)
	m.tableView.SetSize(mainW, bodyH)
	m.queryBar.SetWidth(m.width)
	editorW := m.width - 10
	if editorW > 100 {
		editorW = 100
	}
	m.editor.Resize(editorW)
}

func (m *Model) toastCmd(text string, isErr bool) tea.Cmd {
	m.toast = toastMsg{
		text:  text,
		until: time.Now().Add(3 * time.Second),
		isErr: isErr,
	}
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearToastMsg{}
	})
}
