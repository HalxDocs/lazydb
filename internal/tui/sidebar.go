package tui

import (
	"fmt"
	"strings"
)

type TableMeta struct {
	Name  string
	Count int
}

type Sidebar struct {
	tables []TableMeta
	cursor int
	width  int
	height int
}

func NewSidebar(width, height int) Sidebar {
	return Sidebar{
		width:  width,
		height: height,
	}
}

func (s *Sidebar) SetTables(tables []TableMeta) {
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
	return s.tables[s.cursor].Name
}

func (s Sidebar) Render() string {
	var b strings.Builder

	b.WriteString(SidebarTitleStyle.Render("TABLES") + "\n\n")

	for i, t := range s.tables {
		count := ""
		if t.Count > 0 {
			count = fmt.Sprintf(" %d", t.Count)
		}
		if i == s.cursor {
			label := fmt.Sprintf("▶ %-14s", truncate(t.Name, 14))
			b.WriteString(SidebarItemActiveStyle.Render(label) +
				SidebarCountActiveStyle.Render(count) + "\n")
		} else {
			label := fmt.Sprintf("  %-14s", truncate(t.Name, 14))
			b.WriteString(SidebarItemStyle.Render(label) +
				SidebarCountStyle.Render(count) + "\n")
		}
	}

	return SidebarStyle.Height(s.height).Render(b.String())
}