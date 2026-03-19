package tui

import (
	"database/sql"
	"fmt"

	"github.com/HalxDocs/lazydb/internal/db"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeQuery
)

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

type tablesLoadedMsg struct{ tables []TableMeta }
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
		sidebarWidth := 28
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
			return m, m.loadRows(msg.tables[0].Name)
		}
		return m, nil

	case rowsLoadedMsg:
		if sqlRows, ok := msg.rows.(rowsResult); ok {
			if realRows, ok := sqlRows.rows.(*sql.Rows); ok {
				if err := m.tableView.Load(realRows); err != nil {
					m.err = err
					return m, nil
				}
				realRows.Close()
			}
		}
		m.status = fmt.Sprintf("%d tables · %d rows loaded", len(m.sidebar.tables), len(m.tableView.rows))
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if m.mode == ModeQuery {
			return m.handleQueryMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	if m.mode == ModeQuery {
		var cmd tea.Cmd
		m.queryBar.input, cmd = m.queryBar.input.Update(msg)
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
		return ErrorStyle.Render(fmt.Sprintf("\n  error: %s\n\n  press esc to dismiss", m.err))
	}

	sidebar := m.sidebar.Render()
	main := m.tableView.Render()
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	statusLeft := StatusBarTextStyle.Render(
		fmt.Sprintf("  lazydb  •  %s  •  %s",
			m.sidebar.SelectedTable(),
			m.status,
		),
	)
	statusRight := StatusBarTextStyle.Render("↑↓ rows  ←→ tables  / query  q quit  ")
	spacer := lipgloss.NewStyle().
		Width(m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight)).
		Render("")
	statusBar := StatusBarStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, statusLeft, spacer, statusRight),
	)

	queryBar := ""
	if m.queryBar.IsVisible() {
		queryBar = m.queryBar.Render()
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, queryBar, statusBar)
}

type rowsResult struct {
	rows interface{ Close() error }
}

func (m Model) loadTables() tea.Cmd {
	return func() tea.Msg {
		names, err := m.db.Tables()
		if err != nil {
			return errMsg{err}
		}
		tables := make([]TableMeta, len(names))
		for i, name := range names {
			count, _ := m.db.CountRows(name)
			tables[i] = TableMeta{Name: name, Count: count}
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