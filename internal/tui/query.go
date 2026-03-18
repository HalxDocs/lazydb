package tui

import "github.com/charmbracelet/bubbles/textinput"

type QueryBar struct {
	input   textinput.Model
	visible bool
}

func NewQueryBar() QueryBar {
	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM users..."
	ti.CharLimit = 500
	ti.Width = 80

	return QueryBar{
		input:   ti,
		visible: false,
	}
}

func (q *QueryBar) Show() {
	q.visible = true
	q.input.Focus()
}

func (q *QueryBar) Hide() {
	q.visible = false
	q.input.Blur()
	q.input.SetValue("")
}

func (q *QueryBar) IsVisible() bool {
	return q.visible
}

func (q *QueryBar) Value() string {
	return q.input.Value()
}

func (q QueryBar) Render() string {
	if !q.visible {
		return ""
	}
	prompt := QueryPromptStyle.Render("▶ ")
	return QueryBarStyle.Render(prompt + q.input.View())
}