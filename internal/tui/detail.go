package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderDetail(columns []string, row Row, width, height int) string {
	if len(columns) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(ModalTitleStyle.Render("row detail") + "\n\n")

	// figure column name width
	nameW := 0
	for _, c := range columns {
		if len(c) > nameW {
			nameW = len(c)
		}
	}
	if nameW > 28 {
		nameW = 28
	}

	contentW := width - 12 - nameW
	if contentW < 20 {
		contentW = 20
	}

	for i, col := range columns {
		val := ""
		kind := kindString
		if i < len(row.Cells) {
			val = row.Cells[i]
			kind = row.Kinds[i]
		}
		valStyle := ModalValueStyle
		switch kind {
		case kindNumeric:
			valStyle = lipgloss.NewStyle().Foreground(colorNumeric).Background(colorSurfaceAlt)
		case kindNull:
			valStyle = lipgloss.NewStyle().Foreground(colorNull).Italic(true).Background(colorSurfaceAlt)
		case kindBool:
			valStyle = lipgloss.NewStyle().Foreground(colorBool).Background(colorSurfaceAlt)
		}

		// wrap long values
		wrapped := wordWrap(val, contentW)

		key := ModalKeyStyle.Render(fmt.Sprintf("%-*s", nameW, truncate(col, nameW)))
		first := true
		for _, line := range strings.Split(wrapped, "\n") {
			if first {
				b.WriteString("  " + key + "  " + valStyle.Render(line) + "\n")
				first = false
			} else {
				b.WriteString("  " + strings.Repeat(" ", nameW) + "  " + valStyle.Render(line) + "\n")
			}
		}
	}

	b.WriteString("\n" + ModalDimStyle.Render("y") + ModalDimStyle.Render(" yank cell  ·  ") +
		ModalDimStyle.Render("esc") + ModalDimStyle.Render(" close"))

	modal := ModalStyle.Render(b.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}

func wordWrap(s string, w int) string {
	if w <= 0 || len(s) <= w {
		return s
	}
	var b strings.Builder
	for len(s) > w {
		b.WriteString(s[:w])
		b.WriteString("\n")
		s = s[w:]
	}
	b.WriteString(s)
	return b.String()
}
