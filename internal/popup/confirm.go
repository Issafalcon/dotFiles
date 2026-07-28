// This file implements the install confirmation dialog popup.
//
// When a user chooses to install a module, this dialog appears showing what
// will be installed and asking for confirmation. Buttons: Yes / No / Review
// (Review opens the install.sh / uninstall.sh viewer when a script exists).
package popup

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ConfirmAction indicates whether the confirm dialog is for install or uninstall.
type ConfirmAction string

const (
	ActionInstall   ConfirmAction = "install"
	ActionReinstall ConfirmAction = "reinstall"
	ActionUninstall ConfirmAction = "uninstall"
)

// ConfirmYesMsg is sent when the user confirms the action.
type ConfirmYesMsg struct {
	ModuleName string
	Action     ConfirmAction
}

// ConfirmNoMsg is sent when the user cancels or declines.
type ConfirmNoMsg struct{}

// ConfirmReviewMsg is sent when the user wants to review the module script.
type ConfirmReviewMsg struct {
	ModuleName string
	Action     ConfirmAction
}

// ConfirmModel is the confirmation dialog.
type ConfirmModel struct {
	Model
	moduleName string
	action     ConfirmAction
	items      []string
	subtitle   string
	hasScript  bool
	confirmed  bool
	cursor     int // 0 = Yes, 1 = No, 2 = Review (if hasScript)
}

// NewConfirmDialog creates a confirmation dialog for installing a module.
func NewConfirmDialog(moduleName string, deps []string, hasScript bool) ConfirmModel {
	return ConfirmModel{
		Model:      NewPopup(fmt.Sprintf("Install %s?", moduleName), "", 54, 16).Show(),
		moduleName: moduleName,
		action:     ActionInstall,
		items:      deps,
		subtitle:   "The following will be installed:",
		hasScript:  hasScript,
		cursor:     0,
	}
}

// NewReinstallDialog confirms forcing a re-run of install.sh for one module.
func NewReinstallDialog(moduleName string, hasScript bool) ConfirmModel {
	return ConfirmModel{
		Model:      NewPopup(fmt.Sprintf("Re-run %s install?", moduleName), "", 54, 14).Show(),
		moduleName: moduleName,
		action:     ActionReinstall,
		items: []string{
			"Re-run install.sh (and stow if enabled)",
			"Ignores installed / satisfied status",
			"Dependencies are not reinstalled",
		},
		subtitle:  "Force reinstall:",
		hasScript: hasScript,
		cursor:    0,
	}
}

// NewUninstallDialog creates a confirmation dialog for uninstalling a module.
func NewUninstallDialog(moduleName string, items []string, hasScript bool) ConfirmModel {
	return ConfirmModel{
		Model:      NewPopup(fmt.Sprintf("Uninstall %s?", moduleName), "", 54, 16).Show(),
		moduleName: moduleName,
		action:     ActionUninstall,
		items:      items,
		subtitle:   "The following will be removed:",
		hasScript:  hasScript,
		cursor:     0,
	}
}

func (m ConfirmModel) maxCursor() int {
	if m.hasScript {
		return 2
	}
	return 1
}

// Update handles keyboard input for the confirmation dialog.
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l", "tab":
			if m.cursor < m.maxCursor() {
				m.cursor++
			}
		case "r":
			if m.hasScript {
				return m, func() tea.Msg {
					return ConfirmReviewMsg{ModuleName: m.moduleName, Action: m.action}
				}
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.confirmed = true
				return m, func() tea.Msg {
					return ConfirmYesMsg{ModuleName: m.moduleName, Action: m.action}
				}
			case 2:
				if m.hasScript {
					return m, func() tea.Msg {
						return ConfirmReviewMsg{ModuleName: m.moduleName, Action: m.action}
					}
				}
				fallthrough
			default:
				return m, func() tea.Msg { return ConfirmNoMsg{} }
			}
		case "esc":
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}
	return m, nil
}

// View builds the inner content of the confirmation dialog.
func (m ConfirmModel) View() string {
	var b strings.Builder
	b.WriteString(theme.Subtitle.Render(m.subtitle))
	b.WriteString("\n\n")
	for _, item := range m.items {
		b.WriteString(theme.NormalText.Render("  • " + item))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	btn := func(label string, active bool, bg color.Color) string {
		s := lipgloss.NewStyle().Padding(0, 2)
		if active {
			s = s.Bold(true).Foreground(theme.ColorBackground).Background(bg)
		} else {
			s = s.Foreground(theme.ColorForegroundDim)
		}
		return s.Render(" " + label + " ")
	}

	parts := []string{
		btn("Yes", m.cursor == 0, theme.ColorGreen),
		btn("No", m.cursor == 1, theme.ColorRed),
	}
	if m.hasScript {
		parts = append(parts, btn("Review", m.cursor == 2, theme.ColorCyan))
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, joinWithSpaces(parts)...))
	if m.hasScript {
		b.WriteString("\n")
		b.WriteString(theme.HelpStyle.Render("r review script  ·  ←/→ select  ·  enter confirm"))
	}
	return b.String()
}

func joinWithSpaces(parts []string) []string {
	out := make([]string, 0, len(parts)*2-1)
	for i, p := range parts {
		if i > 0 {
			out = append(out, "  ")
		}
		out = append(out, p)
	}
	return out
}

// Render produces the bordered, centered confirmation dialog.
func (m ConfirmModel) Render(screenWidth, screenHeight int) string {
	if !m.Model.IsVisible() {
		return ""
	}
	return renderPopup(
		m.Model.title, m.View(),
		m.Model.width, m.Model.height,
		screenWidth, screenHeight,
	)
}
