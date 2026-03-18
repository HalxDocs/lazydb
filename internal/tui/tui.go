package tui

import (
	"fmt"

	"github.com/HalxDocs/lazydb/internal/db"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mode represents what the user is currently doing
type Mode int

const (
	ModeNormal Mode = iota
	ModeQuery
)

// Model is the root bubbletea model
type Model struct {
	db        db.DB
	sidebar   Sidebar
	tableView TableView
	queryBar  QueryBar
	mode      Mode
	width     int
	height    int
	err       error
	status    string
	ready     bool
}

// messages
type tablesLoadedMsg struct{ tables []string }
type rowsLoadedMsg struct{ rows interface{} }
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func NewModel(database db.DB) Model {
	return Model{
		db:       database,
		queryBar: NewQueryBar(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadTables()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		sidebarWidth := 24
		mainWidth := m.width - sidebarWidth
		mainHeight := m.height - 3
		m.sidebar = NewSidebar(sidebarWidth, mainHeight)
		m.tableView = NewTableView(mainWidth, mainHeight)
		m.ready = true
		return m, m.loadTables()

	case tablesLoadedMsg:
		m.sidebar.SetTables(msg.tables)
		m.status = fmt.Sprintf("%d tables", len(msg.tables))
		if len(msg.tables) > 0 {
			return m, m.loadRows(msg.tables[0])
		}
		return m, nil

	case rowsLoadedMsg:
		if sqlRows, ok := msg.rows.(rowsResult); ok {
			if err := m.tableView.Load(sqlRows.rows); err != nil {
				m.err = err
			}
			sqlRows.rows.Close()
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		// query mode — route keys to input
		if m.mode == ModeQuery {
			return m.handleQueryMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	// update query bar input when in query mode
	if m.mode == ModeQuery {
		var cmd tea.Cmd
		updated := m.queryBar.input
		updated, cmd = updated.Update(msg)
		m.queryBar.input = updated
		return m, cmd
	}

	return m, nil
}

func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		m.tableView.MoveUp()

	case "down", "j":
		m.tableView.MoveDown()

	case "left", "h":
		m.sidebar.MoveUp()
		table := m.sidebar.SelectedTable()
		if table != "" {
			return m, m.loadRows(table)
		}

	case "right", "l":
		m.sidebar.MoveDown()
		table := m.sidebar.SelectedTable()
		if table != "" {
			return m, m.loadRows(table)
		}

	case "/":
		m.mode = ModeQuery
		m.queryBar.Show()

	case "esc":
		m.err = nil
	}

	return m, nil
}

func (m Model) handleQueryMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.queryBar.Hide()
		return m, nil

	case "enter":
		query := m.queryBar.Value()
		m.mode = ModeNormal
		m.queryBar.Hide()
		if query != "" {
			return m, m.runQuery(query)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.queryBar.input, cmd = m.queryBar.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "\n  connecting to database..."
	}

	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("\n  error: %s\n  press esc to dismiss", m.err))
	}

	// sidebar
	sidebar := m.sidebar.Render()

	// main content
	main := m.tableView.Render()

	// join sidebar and main horizontally
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	// status bar
	statusLeft := StatusBarTextStyle.Render(
		fmt.Sprintf("  lazydb  •  %s  •  %s",
			m.sidebar.SelectedTable(),
			m.status,
		),
	)
	statusRight := StatusBarTextStyle.Render("↑↓ rows  ←→ tables  / query  q quit  ")
	statusBar := StatusBarStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, statusLeft,
			lipgloss.NewStyle().Width(m.width-lipgloss.Width(statusLeft)-lipgloss.Width(statusRight)).Render(""),
			statusRight,
		),
	)

	// query bar
	queryBar := ""
	if m.queryBar.IsVisible() {
		queryBar = m.queryBar.Render()
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, queryBar, statusBar)
}

// commands

type rowsResult struct {
	rows interface{ Close() error }
}

func (m Model) loadTables() tea.Cmd {
	return func() tea.Msg {
		tables, err := m.db.Tables()
		if err != nil {
			return errMsg{err}
		}
		return tablesLoadedMsg{tables}
	}
}

func (m Model) loadRows(table string) tea.Cmd {
	return func() tea.Msg {
		rows, err := m.db.Rows(table, 100)
		if err != nil {
			return errMsg{err}
		}
		return rowsLoadedMsg{rowsResult{rows}}
	}
}

func (m Model) runQuery(query string) tea.Cmd {
	return func() tea.Msg {
		rows, err := m.db.Query(query)
		if err != nil {
			return errMsg{err}
		}
		return rowsLoadedMsg{rowsResult{rows}}
	}
}