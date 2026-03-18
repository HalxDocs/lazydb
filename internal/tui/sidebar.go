package tui

import (
	"fmt"
	"strings"
)

type Sidebar struct {
	tables  []string
	cursor  int
	width   int
	height  int
}

func NewSidebar(width, height int) Sidebar {
	return Sidebar{
		width:  width,
		height: height,
	}
}

func (s *Sidebar) SetTables(tables []string) {
	s.tables = tables
	s.cursor = 0
}

func (s *Sidebar) MoveUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *Sidebar) MoveDown() {
	if s.cursor < len(s.tables)-1 {
		s.cursor++
	}
}

func (s *Sidebar) SelectedTable() string {
	if len(s.tables) == 0 {
		return ""
	}
	return s.tables[s.cursor]
}

func (s Sidebar) Render() string {
	var b strings.Builder

	b.WriteString(SidebarTitleStyle.Render("TABLES") + "\n")

	for i, table := range s.tables {
		item := fmt.Sprintf(" %s", table)
		if i == s.cursor {
			b.WriteString(SidebarItemActiveStyle.Render("▶ "+table) + "\n")
		} else {
			b.WriteString(SidebarItemStyle.Render(item) + "\n")
		}
	}

	return SidebarStyle.Height(s.height).Render(b.String())
}