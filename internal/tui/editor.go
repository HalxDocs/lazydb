package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

type SQLEditor struct {
	area    textarea.Model
	visible bool
}

func NewSQLEditor() SQLEditor {
	ta := textarea.New()
	ta.Placeholder = "SELECT * FROM users WHERE created_at > now() - interval '1 day';"
	ta.Prompt = "│ "
	ta.CharLimit = 5000
	ta.SetWidth(80)
	ta.SetHeight(8)
	ta.ShowLineNumbers = true

	// theme
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		Background(colorSurfaceAlt)
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().
		Foreground(colorAccent).
		Background(colorSurfaceAlt)
	ta.FocusedStyle.Text = lipgloss.NewStyle().
		Foreground(colorText).
		Background(colorSurfaceAlt)
	ta.FocusedStyle.LineNumber = lipgloss.NewStyle().
		Foreground(colorMuted).
		Background(colorSurfaceAlt)
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Background(colorSurfaceAlt)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("#141E30"))

	ta.BlurredStyle = ta.FocusedStyle

	return SQLEditor{area: ta}
}

func (e *SQLEditor) Show() {
	e.visible = true
	e.area.Focus()
}

func (e *SQLEditor) Hide() {
	e.visible = false
	e.area.Blur()
	e.area.Reset()
}

func (e *SQLEditor) IsVisible() bool { return e.visible }
func (e *SQLEditor) Value() string   { return e.area.Value() }

func (e *SQLEditor) Resize(w int) {
	if w < 40 {
		w = 40
	}
	e.area.SetWidth(w - 4)
}

func (e *SQLEditor) SetValue(s string) {
	e.area.SetValue(s)
}

func (e SQLEditor) Render(width, height int) string {
	if !e.visible {
		return ""
	}
	title := ModalTitleStyle.Render("SQL editor")
	hint := ModalDimStyle.Render("Ctrl+Enter to run  ·  Esc to close")
	body := strings.TrimRight(e.area.View(), "\n")
	content := title + "\n\n" + body + "\n\n" + hint
	modal := ModalStyle.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}
