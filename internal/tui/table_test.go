package tui

import (
	"strings"
	"testing"
	"time"
)

func TestFormatValue_Kinds(t *testing.T) {
	cases := []struct {
		in   any
		want cellKind
	}{
		{nil, kindNull},
		{true, kindBool},
		{int64(42), kindNumeric},
		{3.14, kindNumeric},
		{"123", kindNumeric},
		{"hello", kindString},
		{[]byte("9"), kindNumeric},
		{[]byte("abc"), kindString},
		{time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC), kindString},
	}
	for _, c := range cases {
		_, got := formatValue(c.in)
		if got != c.want {
			t.Errorf("formatValue(%v): kind=%v, want=%v", c.in, got, c.want)
		}
	}
}

func TestSortAndFilter(t *testing.T) {
	tv := NewTableView(120, 30)
	tv.columns = []string{"id", "name"}
	tv.allRows = []Row{
		{Cells: []string{"3", "carol"}, Kinds: []cellKind{kindNumeric, kindString}},
		{Cells: []string{"1", "alice"}, Kinds: []cellKind{kindNumeric, kindString}},
		{Cells: []string{"2", "bob"}, Kinds: []cellKind{kindNumeric, kindString}},
	}
	tv.applyFilterSort()

	if tv.RowCount() != 3 {
		t.Fatalf("expected 3 rows, got %d", tv.RowCount())
	}

	tv.colCursor = 0
	tv.CycleSort() // asc on id
	if tv.rows[0].Cells[0] != "1" {
		t.Fatalf("asc sort failed, first row id=%q", tv.rows[0].Cells[0])
	}
	tv.CycleSort() // desc
	if tv.rows[0].Cells[0] != "3" {
		t.Fatalf("desc sort failed, first row id=%q", tv.rows[0].Cells[0])
	}

	tv.SetFilter("ali")
	if tv.RowCount() != 1 || tv.rows[0].Cells[1] != "alice" {
		t.Fatalf("filter failed: rows=%d", tv.RowCount())
	}
}

func TestTruncateAndHumanCount(t *testing.T) {
	if got := truncate("hello world", 5); got != "hell…" {
		t.Errorf("truncate: got %q", got)
	}
	if got := humanCount(1500); got != "1.5k" {
		t.Errorf("humanCount(1500)=%q", got)
	}
	if got := humanCount(2_500_000); got != "2.5M" {
		t.Errorf("humanCount(2.5M)=%q", got)
	}
	if !strings.Contains(humanCount(42), "42") {
		t.Errorf("humanCount(42)")
	}
}
