// This file implements a scrollable read-only viewer for module install/uninstall scripts.
//
// Used from the confirm dialog when the user presses 'r' (Review) to inspect
// install.sh or uninstall.sh before confirming.
package popup

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ScriptDismissMsg is sent when the user closes the script viewer (esc).
// The parent should return to the confirm dialog.
type ScriptDismissMsg struct{}

// ScriptModel is a scrollable popup showing the contents of a shell script.
type ScriptModel struct {
	Model
	lines  []string
	offset int // first visible line index
}

// NewScriptViewer creates a script viewer popup with the given title and content.
func NewScriptViewer(title, content string) ScriptModel {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	return ScriptModel{
		Model:  NewPopup(title, "", 72, 22).Show(),
		lines:  lines,
		offset: 0,
	}
}

// Update handles scrolling and dismiss for the script viewer.
func (m ScriptModel) Update(msg tea.Msg) (ScriptModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		visible := m.visibleLines()
		maxOffset := len(m.lines) - visible
		if maxOffset < 0 {
			maxOffset = 0
		}
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return ScriptDismissMsg{} }
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			if m.offset < maxOffset {
				m.offset++
			}
		case "pgup", "ctrl+u":
			m.offset -= visible
			if m.offset < 0 {
				m.offset = 0
			}
		case "pgdown", "ctrl+d":
			m.offset += visible
			if m.offset > maxOffset {
				m.offset = maxOffset
			}
		case "home", "g":
			m.offset = 0
		case "end", "G":
			m.offset = maxOffset
		}
	}
	return m, nil
}

func (m ScriptModel) visibleLines() int {
	// Inner height minus footer hint line.
	h := m.Model.height - 3
	if h < 3 {
		h = 3
	}
	return h
}

// View renders the visible slice of script lines plus a footer hint.
func (m ScriptModel) View() string {
	var b strings.Builder
	visible := m.visibleLines()
	end := m.offset + visible
	if end > len(m.lines) {
		end = len(m.lines)
	}
	for _, line := range m.lines[m.offset:end] {
		b.WriteString(theme.DimText.Render(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(theme.HelpStyle.Render("j/k scroll  ·  esc back to confirm"))
	return b.String()
}

// Render produces the bordered, centered script viewer.
func (m ScriptModel) Render(screenWidth, screenHeight int) string {
	if !m.Model.IsVisible() {
		return ""
	}
	return renderPopup(
		m.Model.title, m.View(),
		m.Model.width, m.Model.height,
		screenWidth, screenHeight,
	)
}
