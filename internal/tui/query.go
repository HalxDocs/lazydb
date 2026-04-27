package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type QueryBar struct {
	input    textinput.Model
	visible  bool
	history  []string
	histIdx  int // -1 = current draft (not browsing)
	draft    string
	prompt   string
	mode     QueryMode
}

type QueryMode int

const (
	QueryModeSQL    QueryMode = iota // running raw SQL
	QueryModeFilter                  // filter rows
)

func NewQueryBar() QueryBar {
	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM users..."
	ti.CharLimit = 1000
	ti.Width = 80
	ti.Prompt = ""

	ti.PromptStyle = lipgloss.NewStyle().Foreground(colorAccent)
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorMuted)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	return QueryBar{
		input:   ti,
		histIdx: -1,
		prompt:  "❯",
		mode:    QueryModeSQL,
	}
}

func (q *QueryBar) ShowSQL() {
	q.mode = QueryModeSQL
	q.prompt = "❯"
	q.input.Placeholder = "SELECT * FROM users..."
	q.visible = true
	q.input.Focus()
	q.histIdx = -1
}

func (q *QueryBar) ShowFilter(initial string) {
	q.mode = QueryModeFilter
	q.prompt = "filter:"
	q.input.Placeholder = "substring across all cells"
	q.input.SetValue(initial)
	q.visible = true
	q.input.Focus()
}

func (q *QueryBar) Hide() {
	q.visible = false
	q.input.Blur()
	q.input.SetValue("")
	q.histIdx = -1
}

func (q *QueryBar) IsVisible() bool { return q.visible }
func (q *QueryBar) Mode() QueryMode { return q.mode }
func (q *QueryBar) Value() string   { return q.input.Value() }
func (q *QueryBar) SetWidth(w int) {
	if w < 20 {
		w = 20
	}
	q.input.Width = w - 12
}

func (q *QueryBar) PushHistory(s string) {
	if s == "" {
		return
	}
	if len(q.history) > 0 && q.history[len(q.history)-1] == s {
		return
	}
	q.history = append(q.history, s)
	if len(q.history) > 50 {
		q.history = q.history[len(q.history)-50:]
	}
}

func (q *QueryBar) HistoryUp() {
	if len(q.history) == 0 {
		return
	}
	if q.histIdx == -1 {
		q.draft = q.input.Value()
		q.histIdx = len(q.history) - 1
	} else if q.histIdx > 0 {
		q.histIdx--
	}
	q.input.SetValue(q.history[q.histIdx])
	q.input.CursorEnd()
}

func (q *QueryBar) HistoryDown() {
	if q.histIdx == -1 {
		return
	}
	q.histIdx++
	if q.histIdx >= len(q.history) {
		q.histIdx = -1
		q.input.SetValue(q.draft)
	} else {
		q.input.SetValue(q.history[q.histIdx])
	}
	q.input.CursorEnd()
}

func (q QueryBar) Render() string {
	if !q.visible {
		return ""
	}
	style := QueryBarStyle
	promptStyle := QueryPromptStyle
	if q.mode == QueryModeFilter {
		style = FilterBarStyle
		promptStyle = FilterPromptStyle
	}
	prompt := promptStyle.Render(" " + q.prompt + " ")
	return style.Render(prompt + q.input.View())
}
