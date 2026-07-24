// Package detail — overview.go implements the Overview tab for the detail panel.
//
// The Overview tab shows information about the currently selected module:
//
//   - Module name and description
//   - Website and repository URLs
//   - A navigable list of external dependencies with their install status
//
// # Dependency Checking
//
// Each dependency has a "checking" state represented by a spinner. When the
// user selects a module, the app runs shell commands to check if dependencies
// are installed (e.g., "rg --version" for ripgrep). Results arrive as
// DepCheckResultMsg messages and update the status indicators.
//
// # Spinner Component
//
// The spinner comes from the Bubbles component library. It's a self-contained
// model that manages its own animation frames via tea.Cmd tick functions.
// We embed it in OverviewModel and forward its messages in Update().
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner
// See: https://pkg.go.dev/charm.land/bubbletea/v2
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

	// Spinner is a Bubbles component that shows an animated loading indicator.
	// It auto-ticks via tea.Cmd, so you just forward its Update() and View().
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner
	"charm.land/bubbles/v2/spinner"

	// Theme provides shared styles and colours for the app.
	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Dependency Status
// ---------------------------------------------------------------------------

// DepStatus represents the current install status of a single external
// dependency. This is a view-specific struct — it mirrors information from
// module.ExternalDep but adds UI state (Installed, Checking) that only the
// detail panel cares about.
//
// In Go, structs are value types. When you assign a struct to another
// variable, Go copies all the fields. This means modifying the copy
// doesn't affect the original — which is important when updating slices.
//
// See: https://go.dev/tour/moretypes/2
// See: https://go.dev/doc/effective_go#composite_literals
type DepStatus struct {
	Name     string // Human-readable dependency name (e.g., "ripgrep").
	Method   string // Install method / package manager (e.g., "apt", "brew", "cargo").
	Installed bool  // Whether the dependency is confirmed installed.
	Checking  bool  // Whether we're currently running the check command.
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// DepCheckResultMsg is sent when an async dependency check completes.
// The installer or command runner sends this message after running a
// dependency's CheckCommand and observing its exit code.
//
// Custom messages are the primary way to communicate async results back
// to a Bubble Tea model. Any Go type can be a message.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg
type DepCheckResultMsg struct {
	Name      string // Which dependency was checked.
	Installed bool   // Whether the check command succeeded (exit 0 = installed).
}

// InstallDepMsg is sent when the user presses Enter on a dependency to
// trigger installation. The parent model can listen for this and kick
// off the actual install process.
type InstallDepMsg struct {
	Name   string // The dependency to install.
	Method string // The package manager to use.
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// OverviewModel holds all state for the Overview tab.
//
// Key Go concepts used here:
//
//   - Slice ([]DepStatus): A dynamically-sized, reference-backed list.
//     See: https://go.dev/blog/slices-intro
//
//   - Struct embedding vs composition: We store spinner.Model as a named
//     field (not embedded) so we can namespace its methods clearly.
//     See: https://go.dev/doc/effective_go#embedding
type OverviewModel struct {
	// Module metadata displayed at the top of the tab.
	moduleName  string
	description string
	website     string
	repo        string

	// deps is the list of external dependencies and their install status.
	// Slices in Go are reference types backed by an array. Appending may
	// allocate a new array if capacity is exceeded.
	// See: https://go.dev/tour/moretypes/7
	deps []DepStatus

	// depCursor tracks which dependency is currently highlighted.
	// The user navigates with j/k keys.
	depCursor int

	// width and height define the rendering area.
	width  int
	height int

	// spinner provides an animated indicator for dependencies being checked.
	// spinner.Model is from the Bubbles component library — it manages its
	// own tick-based animation. We forward messages to it in Update().
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner
	spinner spinner.Model
}

// NewOverviewModel creates an initialised OverviewModel.
//
// We configure the spinner with the "Dot" style (a sequence of braille dots
// that animate smoothly). The spinner's style is set to yellow to match our
// "checking" status colour.
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner#New
func NewOverviewModel(width, height int) OverviewModel {
	// Create a new spinner with the Dot animation pattern.
	// spinner.New() returns a spinner.Model value.
	s := spinner.New()

	// spinner.Dot is a predefined animation: ⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷
	// Other options include spinner.Line, spinner.MiniDot, spinner.Globe, etc.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner#Spinner
	s.Spinner = spinner.Dot

	// Style the spinner text with yellow to indicate "in progress".
	s.Style = lipgloss.NewStyle().Foreground(theme.ColorYellow)

	return OverviewModel{
		width:   width,
		height:  height,
		spinner: s,
	}
}

// SetModule updates the overview to display a new module's information.
// This is called by the parent when the user selects a different module
// in the sidebar.
//
// In Go, methods with pointer receivers (*OverviewModel) can modify the
// receiver. Value receivers (OverviewModel) work on a copy and changes
// are discarded. Use pointer receivers when you need to mutate state.
//
// See: https://go.dev/tour/methods/4
// See: https://go.dev/doc/effective_go#pointers_vs_values
func (m *OverviewModel) SetModule(name, description, website, repo string, deps []DepStatus) {
	m.moduleName = name
	m.description = description
	m.website = website
	m.repo = repo
	m.deps = deps
	// Reset cursor to the top of the dependency list when switching modules.
	m.depCursor = 0
}

// SetSize updates the rendering dimensions. Called on terminal resize.
func (m *OverviewModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages for the Overview tab.
//
// # Message Flow
//
// 1. Key presses (j/k/enter) navigate the dependency list or trigger installs.
// 2. DepCheckResultMsg updates a dependency's installed/checking state.
// 3. spinner.TickMsg is forwarded to the spinner for animation.
//
// Notice we return (OverviewModel, tea.Cmd) — the concrete type, not tea.Model.
// This is the idiomatic pattern for sub-models that are managed by a parent.
//
// See: https://go.dev/doc/effective_go#type_switch
func (m OverviewModel) Update(msg tea.Msg) (OverviewModel, tea.Cmd) {
	// Forward spinner messages first. The spinner generates its own tick
	// commands to animate. We must always forward these or the animation stops.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner#Model.Update
	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {

	// Handle keyboard input for dependency list navigation.
	case tea.KeyPressMsg:
		switch msg.String() {

		// "j" or "down" moves the cursor down in the dependency list.
		case "j", "down":
			if len(m.deps) > 0 {
				// Clamp to the last item. min() is not built-in for ints in
				// older Go, but Go 1.21+ has it. We do it manually for clarity.
				m.depCursor++
				if m.depCursor >= len(m.deps) {
					m.depCursor = len(m.deps) - 1
				}
			}
			return m, spinnerCmd

		// "k" or "up" moves the cursor up.
		case "k", "up":
			if m.depCursor > 0 {
				m.depCursor--
			}
			return m, spinnerCmd

		// "enter" triggers installation of the highlighted dependency.
		case "enter":
			if len(m.deps) > 0 && m.depCursor < len(m.deps) {
				dep := m.deps[m.depCursor]
				// Only trigger install if the dependency is NOT already installed.
				if !dep.Installed && !dep.Checking {
					// Return a command that produces an InstallDepMsg.
					// The parent model will pick this up and start the install.
					installCmd := func() tea.Msg {
						return InstallDepMsg{
							Name:   dep.Name,
							Method: dep.Method,
						}
					}
					// tea.Batch combines multiple commands to run concurrently.
					// See: https://pkg.go.dev/charm.land/bubbletea/v2#Batch
					return m, tea.Batch(spinnerCmd, installCmd)
				}
			}
			return m, spinnerCmd
		}

	// DepCheckResultMsg arrives when an async dependency check finishes.
	// We find the matching dependency by name and update its status.
	case DepCheckResultMsg:
		for i := range m.deps {
			if m.deps[i].Name == msg.Name {
				m.deps[i].Installed = msg.Installed
				m.deps[i].Checking = false
			}
		}
		return m, spinnerCmd
	}

	return m, spinnerCmd
}

// View renders the Overview tab content.
//
// The layout is:
//
//	┌─────────────────────────────────────┐
//	│ 📦 Module Name                      │
//	│ Description text                    │
//	│                                     │
//	│ 🌐 Website: https://...            │
//	│ 📂 Repo:    https://github.com/... │
//	│                                     │
//	│ Dependencies:                       │
//	│  ✓ ripgrep        (apt)            │
//	│  ✗ fd-find        (apt)   ← cursor │
//	│  ⟳ bat            (brew)           │
//	│                                     │
//	│ Press Enter to install              │
//	└─────────────────────────────────────┘
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2#Style.Render
func (m OverviewModel) View() string {
	// If no module is selected yet, show a placeholder message.
	if m.moduleName == "" {
		// lipgloss.Place centres the message within the available space.
		// See: https://pkg.go.dev/charm.land/lipgloss/v2#Place
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			theme.DimText.Render("Select a module from the sidebar"),
		)
	}

	// Use a strings.Builder for efficient string concatenation.
	// strings.Builder minimises memory allocations when building strings
	// incrementally, unlike += which creates a new string each time.
	// See: https://pkg.go.dev/strings#Builder
	var b strings.Builder

	// --- Module Header ---
	// Render the module name as a bold pink title.
	icon := theme.GetModuleIcon(m.moduleName)
	b.WriteString(theme.Title.Render(fmt.Sprintf("%s %s", icon, m.moduleName)))
	b.WriteString("\n")

	// Render the description in normal text.
	if m.description != "" {
		b.WriteString(theme.NormalText.Render(m.description))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// --- URLs ---
	// Render website and repo links with the cyan underlined URL style.
	if m.website != "" {
		label := theme.Subtitle.Render("Website: ")
		url := theme.URLStyle.Render(m.website)
		b.WriteString(label + url + "\n")
	}
	if m.repo != "" {
		label := theme.Subtitle.Render("Repo:    ")
		url := theme.URLStyle.Render(m.repo)
		b.WriteString(label + url + "\n")
	}
	b.WriteString("\n")

	// --- Dependency Table ---
	if len(m.deps) > 0 {
		b.WriteString(theme.Subtitle.Render("Dependencies:"))
		b.WriteString("\n\n")

		// Iterate over dependencies and render each row.
		// The range keyword iterates over slices, returning (index, value).
		// See: https://go.dev/tour/moretypes/16
		for i, dep := range m.deps {
			// Determine the status indicator.
			var statusIcon string
			switch {
			case dep.Checking:
				// Show the animated spinner for deps currently being checked.
				statusIcon = m.spinner.View()
			case dep.Installed:
				// theme.StatusInstalled has SetString("✓") with green colour.
				statusIcon = theme.StatusInstalled.String()
			default:
				// theme.StatusNotInstalled has SetString("✗") with red colour.
				statusIcon = theme.StatusNotInstalled.String()
			}

			// Format the dependency name and install method.
			name := dep.Name
			method := theme.DimText.Render(fmt.Sprintf("(%s)", dep.Method))

			// Build the row. If this row is the cursor position, highlight it.
			row := fmt.Sprintf("  %s  %-20s %s", statusIcon, name, method)

			if i == m.depCursor {
				// Highlight the selected row with a cyan foreground and a
				// pointer arrow to indicate focus.
				row = lipgloss.NewStyle().
					Foreground(theme.ColorCyan).
					Bold(true).
					Render(fmt.Sprintf("▸ %s  %-20s %s", statusIcon, name, method))
			}

			b.WriteString(row)
			b.WriteString("\n")
		}

		// --- Help Hint ---
		b.WriteString("\n")
		hint := theme.KeyStyle.Render("Enter") +
			theme.DescStyle.Render(" install selected dependency")
		b.WriteString(hint)
	} else {
		// No dependencies — show a message.
		b.WriteString(theme.DimText.Render("No external dependencies for this module."))
	}

	return b.String()
}
