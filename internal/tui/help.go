package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpEntry struct {
	keys string
	desc string
}

var helpSections = []struct {
	title   string
	entries []helpEntry
}{
	{
		"Navigation",
		[]helpEntry{
			{"↑ ↓ k j", "move row up/down"},
			{"← → h l", "move column cursor"},
			{"g / G", "jump to first / last row"},
			{"Ctrl+d / Ctrl+u", "page down / up"},
			{"Tab / Shift+Tab", "next / prev table"},
			{"[ / ]", "next / prev table"},
		},
	},
	{
		"Tables & data",
		[]helpEntry{
			{"Enter", "open row detail"},
			{"s", "toggle schema view"},
			{"o", "sort by current column (asc → desc → none)"},
			{"f", "filter rows"},
			{"r", "refresh current table"},
			{"y", "yank current cell to clipboard"},
			{"Y", "yank current row as TSV"},
			{"e", "export visible rows to CSV"},
		},
	},
	{
		"Sidebar",
		[]helpEntry{
			{"t", "search tables"},
		},
	},
	{
		"Query",
		[]helpEntry{
			{"/", "single-line query"},
			{":", "multi-line SQL editor"},
			{"↑ / ↓ in query", "previous / next history"},
			{"Enter", "run query"},
			{"Ctrl+Enter (editor)", "run query"},
			{"Esc", "close / dismiss"},
		},
	},
	{
		"App",
		[]helpEntry{
			{"?", "toggle this help"},
			{"q / Ctrl+c", "quit"},
		},
	},
}

func renderHelp(width, height int) string {
	var b strings.Builder
	b.WriteString(ModalTitleStyle.Render("lazydb · keyboard shortcuts") + "\n\n")

	for _, sec := range helpSections {
		b.WriteString(ModalTitleStyle.Render(sec.title) + "\n")
		for _, e := range sec.entries {
			line := KbdStyle.Render(" "+e.keys+" ") + "  " +
				ModalValueStyle.Render(e.desc)
			b.WriteString("  " + line + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(ModalDimStyle.Render("press ? or esc to close"))

	modal := ModalStyle.Render(b.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}
