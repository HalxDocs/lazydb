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
	tables   []TableMeta
	filtered []int // indices into tables, in display order
	cursor   int   // index into filtered
	width    int
	height   int
	search   string
	driver   string
}

func NewSidebar(width, height int) Sidebar {
	return Sidebar{
		width:  width,
		height: height,
	}
}

func (s *Sidebar) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *Sidebar) SetDriver(d string) { s.driver = d }

func (s *Sidebar) SetTables(tables []TableMeta) {
	s.tables = tables
	s.cursor = 0
	s.applyFilter()
}

func (s *Sidebar) UpdateCount(name string, count int) {
	for i := range s.tables {
		if s.tables[i].Name == name {
			s.tables[i].Count = count
			return
		}
	}
}

func (s *Sidebar) SetSearch(q string) {
	s.search = q
	s.cursor = 0
	s.applyFilter()
}

func (s *Sidebar) Search() string { return s.search }

func (s *Sidebar) applyFilter() {
	s.filtered = s.filtered[:0]
	if s.search == "" {
		for i := range s.tables {
			s.filtered = append(s.filtered, i)
		}
		return
	}
	q := strings.ToLower(s.search)
	for i, t := range s.tables {
		if strings.Contains(strings.ToLower(t.Name), q) {
			s.filtered = append(s.filtered, i)
		}
	}
}

func (s *Sidebar) MoveUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *Sidebar) MoveDown() {
	if s.cursor < len(s.filtered)-1 {
		s.cursor++
	}
}

func (s *Sidebar) Top()    { s.cursor = 0 }
func (s *Sidebar) Bottom() { s.cursor = max(0, len(s.filtered)-1) }

func (s *Sidebar) SelectedTable() string {
	if len(s.filtered) == 0 {
		return ""
	}
	return s.tables[s.filtered[s.cursor]].Name
}

func (s *Sidebar) SelectedMeta() (TableMeta, bool) {
	if len(s.filtered) == 0 {
		return TableMeta{}, false
	}
	return s.tables[s.filtered[s.cursor]], true
}

func (s *Sidebar) Tables() []TableMeta { return s.tables }

func (s Sidebar) Render() string {
	var b strings.Builder

	driver := s.driver
	if driver == "" {
		driver = "db"
	}
	b.WriteString(SidebarTitleStyle.Render("◆ "+strings.ToUpper(driver)) + "\n")
	b.WriteString(SidebarSubtitleStyle.Render(fmt.Sprintf("%d tables", len(s.tables))) + "\n")

	if s.search != "" {
		b.WriteString("\n")
		b.WriteString(SidebarSearchStyle.Render("/ "+s.search) + "\n")
	}
	b.WriteString("\n")

	visible := s.height - 6
	if s.search != "" {
		visible -= 2
	}
	if visible < 1 {
		visible = 1
	}

	start := 0
	if s.cursor >= visible {
		start = s.cursor - visible + 1
	}
	stop := start + visible
	if stop > len(s.filtered) {
		stop = len(s.filtered)
	}

	nameW := s.width - 8
	if nameW < 8 {
		nameW = 8
	}

	for i := start; i < stop; i++ {
		t := s.tables[s.filtered[i]]
		count := ""
		if t.Count > 0 {
			count = humanCount(t.Count)
		}
		if i == s.cursor {
			label := fmt.Sprintf("▸ %-*s", nameW, truncate(t.Name, nameW))
			b.WriteString(SidebarItemActiveStyle.Render(label) +
				SidebarCountActiveStyle.Render(fmt.Sprintf(" %s", count)) + "\n")
		} else {
			label := fmt.Sprintf("  %-*s", nameW, truncate(t.Name, nameW))
			b.WriteString(SidebarItemStyle.Render(label) +
				SidebarCountStyle.Render(fmt.Sprintf(" %s", count)) + "\n")
		}
	}

	if len(s.filtered) == 0 && s.search != "" {
		b.WriteString(SidebarSubtitleStyle.Render("  no matches") + "\n")
	}

	return SidebarStyle.Width(s.width).Height(s.height).Render(b.String())
}

// humanCount formats integer counts as 1.2k, 3.4M, etc.
func humanCount(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
