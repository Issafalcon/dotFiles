package sidebar

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestSearchEscKeepsFilter(t *testing.T) {
	items := []ModuleItem{
		{Name: "nvim", Description: "editor"},
		{Name: "node", Description: "runtime"},
		{Name: "git", Description: "vcs"},
	}
	m := NewModel(items, 40, 30)

	m, _ = m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	if !m.IsSearching() {
		t.Fatal("expected search mode after /")
	}

	for _, ch := range "nv" {
		m, _ = m.Update(tea.KeyPressMsg{Code: rune(ch), Text: string(ch)})
	}
	if got := len(m.filtered); got != 1 || m.filtered[0].Name != "nvim" {
		t.Fatalf("filter=%v", names(m.filtered))
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.IsSearching() {
		t.Fatal("esc should leave search input")
	}
	if m.searchInput.Value() != "nv" {
		t.Fatalf("filter cleared unexpectedly: %q", m.searchInput.Value())
	}
	if len(m.filtered) != 1 {
		t.Fatalf("filtered list lost after esc: %v", names(m.filtered))
	}

	// Navigate filtered list in normal mode.
	m, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.Selected() != "nvim" {
		t.Fatalf("selected=%q", m.Selected())
	}

	// Second esc clears the filter.
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.searchInput.Value() != "" {
		t.Fatalf("second esc should clear filter, got %q", m.searchInput.Value())
	}
	if len(m.filtered) != 3 {
		t.Fatalf("want full list, got %v", names(m.filtered))
	}
}

func names(items []ModuleItem) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Name
	}
	return out
}
