// Package detail — config.go implements the Configuration tab for the detail panel.
//
// The Configuration tab allows users to view and modify module-specific
// settings before installation. For example, a "dotnet" module might offer
// a choice of .NET SDK version, and a "zsh" module might let you pick a
// plugin manager.
//
// # Config Options
//
// Each option has a name, description, default value, list of valid choices,
// and the user's current selection. The user navigates with j/k and toggles
// selections with Enter.
//
// This is a relatively simple model — it's a navigable list of options with
// no async operations or external components.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2
// See: https://pkg.go.dev/charm.land/lipgloss/v2
package detail

import (
	"fmt"
	"strings"

	// Bubble Tea — the TUI framework.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// Lip Gloss — terminal styling.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	// Theme provides shared styles and colours for the app.
	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Config Option
// ---------------------------------------------------------------------------

// ConfigOption represents a user-configurable setting for a module.
//
// This is a view-specific struct. It mirrors module.ConfigOption but adds a
// Selected field to track the user's current choice within the TUI. The
// module package defines the "data" version; this struct adds "UI state".
//
// # Why a Separate Struct?
//
// Separating data from UI state is a common pattern. The data model
// (module.ConfigOption) stays pure and can be serialised/deserialised
// without UI concerns. The view model (this ConfigOption) adds ephemeral
// state that only matters during the TUI session.
//
// See: https://go.dev/doc/effective_go#composite_literals
type ConfigOption struct {
	Name        string   // Short identifier (e.g., "dotnet_version").
	Description string   // Human-readable label shown to the user.
	Default     string   // Default value if the user doesn't choose.
	Choices     []string // Valid values the user can pick from (empty = freeform).
	Selected    string   // The user's current selection (starts as Default).
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// ConfigChangedMsg is sent when the user changes a configuration option.
// The parent model can listen for this to persist the selection.
type ConfigChangedMsg struct {
	ModuleName string // Which module the config belongs to.
	OptionName string // Which option was changed.
	NewValue   string // The newly selected value.
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// ConfigModel holds all state for the Configuration tab.
//
// It's a simple navigable list: the user moves a cursor (j/k) and presses
// Enter to cycle through available choices for each option.
//
// See: https://go.dev/tour/moretypes/2
type ConfigModel struct {
	// moduleName identifies which module's config is displayed.
	moduleName string

	// options is the list of configurable settings for the current module.
	options []ConfigOption

	// cursor tracks which option is currently highlighted.
	cursor int

	// width and height define the rendering area.
	width  int
	height int
}

// NewConfigModel creates an initialised ConfigModel.
//
// The model starts empty — no module is selected yet. The parent calls
// SetModule() when the user selects a module in the sidebar.
//
// See: https://go.dev/doc/effective_go#composite_literals
func NewConfigModel(width, height int) ConfigModel {
	return ConfigModel{
		width:  width,
		height: height,
	}
}

// SetModule loads configuration options for the given module.
// This replaces the current options and resets the cursor to the top.
//
// Pointer receiver because we're mutating the model.
// See: https://go.dev/tour/methods/4
func (m *ConfigModel) SetModule(name string, options []ConfigOption) {
	m.moduleName = name
	m.options = options
	m.cursor = 0
}

// SetSize updates the rendering dimensions. Called on terminal resize.
func (m *ConfigModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages for the Configuration tab.
//
// Key bindings:
//   - j / down: Move cursor down
//   - k / up: Move cursor up
//   - enter: Cycle to the next choice for the selected option
//
// This is a simple, synchronous model — no async operations or sub-components.
//
// See: https://go.dev/doc/effective_go#type_switch
func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		// Navigate down through the option list.
		case "j", "down":
			if len(m.options) > 0 {
				m.cursor++
				if m.cursor >= len(m.options) {
					m.cursor = len(m.options) - 1
				}
			}
			return m, nil

		// Navigate up through the option list.
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		// Cycle to the next choice for the selected option.
		case "enter":
			if len(m.options) > 0 && m.cursor < len(m.options) {
				opt := &m.options[m.cursor]

				// Only cycle if there are choices to cycle through.
				// If Choices is empty, the option accepts freeform input
				// and cycling doesn't apply.
				if len(opt.Choices) > 0 {
					// Find the index of the current selection in the choices slice.
					// Go doesn't have a built-in indexOf for slices, so we loop.
					// See: https://go.dev/tour/moretypes/16
					currentIdx := -1
					for i, choice := range opt.Choices {
						if choice == opt.Selected {
							currentIdx = i
							break // break exits the innermost for loop.
						}
					}

					// Cycle to the next choice using modular arithmetic.
					// If the current selection isn't found (-1), we start at 0.
					nextIdx := (currentIdx + 1) % len(opt.Choices)
					opt.Selected = opt.Choices[nextIdx]

					// Notify the parent about the change.
					cmd := func() tea.Msg {
						return ConfigChangedMsg{
							ModuleName: m.moduleName,
							OptionName: opt.Name,
							NewValue:   opt.Selected,
						}
					}
					return m, cmd
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the Configuration tab content.
//
// Layout with options:
//
//	┌───────────────────────────────────────┐
//	│ ⚙ Configuration: nvim               │
//	│                                       │
//	│ ▸ Plugin Manager        [lazy.nvim]  │
//	│   Language Servers      [all]         │
//	│   Tree-sitter Parsers   [all]         │
//	│                                       │
//	│ Enter to cycle options                │
//	└───────────────────────────────────────┘
//
// Layout with no options:
//
//	┌───────────────────────────────────────┐
//	│                                       │
//	│   No configuration options for        │
//	│   this module.                        │
//	│                                       │
//	└───────────────────────────────────────┘
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2#Style.Render
func (m ConfigModel) View() string {
	// If no module is selected, show a placeholder.
	if m.moduleName == "" {
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			theme.DimText.Render("Select a module from the sidebar"),
		)
	}

	// If the module has no configuration options, show a centred message.
	if len(m.options) == 0 {
		msg := theme.DimText.Render(
			fmt.Sprintf("No configuration options for %s.", m.moduleName),
		)
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			msg,
		)
	}

	// Build the config list view.
	// strings.Builder provides efficient incremental string construction.
	// See: https://pkg.go.dev/strings#Builder
	var b strings.Builder

	// Header with the module name.
	b.WriteString(theme.Subtitle.Render(
		fmt.Sprintf("%s Configuration: %s", theme.IconGear, m.moduleName),
	))
	b.WriteString("\n\n")

	// Render each config option.
	for i, opt := range m.options {
		// Determine the cursor indicator.
		// The active row gets a cyan arrow (▸), inactive rows get a space.
		var cursor string
		if i == m.cursor {
			cursor = lipgloss.NewStyle().
				Foreground(theme.ColorCyan).
				Bold(true).
				Render("▸ ")
		} else {
			cursor = "  "
		}

		// Render the option name.
		name := theme.NormalText.Render(opt.Name)

		// Render the current selection in brackets.
		// If no selection has been made yet, show the default value.
		displayValue := opt.Selected
		if displayValue == "" {
			displayValue = opt.Default
		}

		// Style the value — use green for selected values, dim for defaults.
		var valueStyle lipgloss.Style
		if opt.Selected != "" && opt.Selected != opt.Default {
			// User has made a custom choice — highlight in green.
			valueStyle = lipgloss.NewStyle().Foreground(theme.ColorGreen)
		} else {
			// Default value — show in dim text.
			valueStyle = lipgloss.NewStyle().Foreground(theme.ColorForegroundDim)
		}
		value := valueStyle.Render(fmt.Sprintf("[%s]", displayValue))

		// Assemble the row.
		row := fmt.Sprintf("%s%-25s %s", cursor, name, value)

		// If this is the selected row, show the description below it.
		if i == m.cursor && opt.Description != "" {
			row += "\n" + "    " + theme.DimText.Render(opt.Description)
		}

		b.WriteString(row)
		b.WriteString("\n")
	}

	// --- Help Hint ---
	b.WriteString("\n")
	hint := theme.KeyStyle.Render("Enter") +
		theme.DescStyle.Render(" cycle options") +
		theme.DimText.Render("  •  ") +
		theme.KeyStyle.Render("j/k") +
		theme.DescStyle.Render(" navigate")
	b.WriteString(hint)

	return b.String()
}
