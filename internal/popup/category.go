// Category picker popup — select All or a module Category to filter the sidebar.
package popup

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// CategorySelectedMsg is sent when the user picks a category ("" means All).
type CategorySelectedMsg struct {
	Category string
}

// CategoryCancelMsg is sent when the user dismisses the picker without changing.
type CategoryCancelMsg struct{}

// CategoryModel is a simple list picker for module categories.
type CategoryModel struct {
	Model
	options []string // first entry is always "All" (maps to "")
	cursor  int
}

// NewCategoryPicker creates a category filter popup.
// categories should be the sorted list from the registry (without "All").
// current is the active filter ("" for All).
func NewCategoryPicker(categories []string, current string) CategoryModel {
	opts := make([]string, 0, len(categories)+1)
	opts = append(opts, "All")
	opts = append(opts, categories...)

	cursor := 0
	if current != "" {
		for i, c := range opts {
			if c == current {
				cursor = i
				break
			}
		}
	}

	return CategoryModel{
		Model:   NewPopup("Filter by category", "", 36, min(12, len(opts)+4)).Show(),
		options: opts,
		cursor:  cursor,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Update handles navigation and selection for the category picker.
func (m CategoryModel) Update(msg tea.Msg) (CategoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			cat := m.options[m.cursor]
			if cat == "All" {
				cat = ""
			}
			return m, func() tea.Msg { return CategorySelectedMsg{Category: cat} }
		case "esc", "q":
			return m, func() tea.Msg { return CategoryCancelMsg{} }
		}
	}
	return m, nil
}

// View renders the category list.
func (m CategoryModel) View() string {
	var b strings.Builder
	b.WriteString(theme.DimText.Render("j/k move  ·  enter select  ·  esc cancel"))
	b.WriteString("\n\n")
	for i, opt := range m.options {
		label := opt
		style := lipgloss.NewStyle().Padding(0, 1)
		if i == m.cursor {
			style = style.Bold(true).
				Foreground(theme.ColorBackground).
				Background(theme.ColorCyan)
			b.WriteString(style.Render("▸ " + label))
		} else {
			b.WriteString(style.Foreground(theme.ColorForegroundDim).Render("  " + label))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Render produces the bordered centered picker.
func (m CategoryModel) Render(screenWidth, screenHeight int) string {
	if !m.Model.IsVisible() {
		return ""
	}
	return renderPopup(
		m.Model.title, m.View(),
		m.Model.width, m.Model.height,
		screenWidth, screenHeight,
	)
}
