// Package detail — output.go implements the Output tab for the detail panel.
//
// The Output tab is a scrollable log viewer that shows the real-time output
// of module installation commands. It uses two Bubbles components:
//
//   - viewport.Model: A scrollable text area that handles its own keyboard
//     input (page up/down, mouse wheel, etc.). Think of it as a "less"
//     pager embedded in the TUI.
//     See: https://pkg.go.dev/charm.land/bubbles/v2/viewport
//
//   - progress.Model: An animated progress bar shown during installations.
//     It smoothly animates between percentage values.
//     See: https://pkg.go.dev/charm.land/bubbles/v2/progress
//
// # Viewport Pattern
//
// The viewport is a "passthrough" component — you forward all messages to it
// via Update(), and it handles scrolling internally. You just need to set its
// content with SetContent() whenever the log lines change.
//
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

	// progress is a Bubbles component that renders animated progress bars.
	// It interpolates smoothly between percentages using tea.Cmd ticks.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/progress
	"charm.land/bubbles/v2/progress"

	// viewport is a Bubbles component that provides a scrollable text area.
	// It handles keyboard/mouse input for scrolling automatically.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport
	"charm.land/bubbles/v2/viewport"

	// Theme provides shared styles and colours for the app.
	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// OutputModel holds all state for the Output (log viewer) tab.
//
// # Viewport Component
//
// viewport.Model is a pre-built scrollable text area from the Bubbles library.
// You set its content with SetContent(string), and it handles scrolling via
// keyboard (page up/down, arrows) and mouse wheel events automatically.
// In Update(), you forward all messages to it so it can process scroll input.
//
// # Progress Component
//
// progress.Model renders an animated progress bar. You set the percentage
// with SetPercent() or render at a specific percentage with ViewAs(). It
// generates its own tick commands for smooth animation.
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport#Model
// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#Model
type OutputModel struct {
	// viewport is the scrollable text area displaying log lines.
	// It's a Bubbles component that manages its own scroll state.
	viewport viewport.Model

	// lines stores all log output lines. We keep them in a slice so we can
	// rebuild the viewport content when new lines are appended.
	//
	// Slices in Go are backed by an underlying array. When you append() and
	// the slice exceeds its capacity, Go allocates a new, larger array and
	// copies the data. The amortised cost is O(1) per append.
	// See: https://go.dev/blog/slices-intro
	lines []string

	// width and height define the rendering area.
	width  int
	height int

	// isInstalling indicates whether a module installation is in progress.
	// When true, we show a progress bar at the top of the tab.
	isInstalling bool

	// currentModule is the name of the module currently being installed.
	currentModule string

	// progress is the animated progress bar shown during installation.
	progress progress.Model
}

// NewOutputModel creates an initialised OutputModel with a viewport and
// progress bar sized to fit the given dimensions.
//
// The viewport is created with the full width and most of the height.
// We reserve 2 lines at the top for the progress bar area when installing.
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport#New
// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#New
func NewOutputModel(width, height int) OutputModel {
	// Create a viewport that fills the available space.
	// In Bubbles v2, viewport.New() uses the functional options pattern:
	// you pass Option functions (WithWidth, WithHeight) instead of positional args.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport#New
	vp := viewport.New(
		viewport.WithWidth(width),
		viewport.WithHeight(height),
	)

	// Create a progress bar. progress.New() also uses functional options.
	//
	// Functional options pattern: instead of a config struct, you pass
	// functions that modify the default configuration. This is a common
	// Go pattern for APIs with many optional settings.
	//
	// WithDefaultBlend() creates a smooth colour gradient for the bar.
	// WithWidth() sets the bar's character width.
	//
	// See: https://go.dev/doc/effective_go#composite_literals
	// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#New
	p := progress.New(
		progress.WithDefaultBlend(),
		progress.WithWidth(width),
	)

	return OutputModel{
		viewport: vp,
		lines:    make([]string, 0),
		width:    width,
		height:   height,
		progress: p,
	}
}

// AppendLine adds a new line to the output log and updates the viewport.
// The viewport is automatically scrolled to the bottom so the user sees
// the latest output — like "tail -f" in a terminal.
//
// This method uses a pointer receiver (*OutputModel) because it modifies
// the model's state (appending to lines, updating viewport content).
//
// See: https://go.dev/tour/methods/4
func (m *OutputModel) AppendLine(line string) {
	// append() is a built-in Go function that adds elements to a slice.
	// It returns a new slice header (which may point to a new backing array
	// if the capacity was exceeded). We must reassign the result.
	// See: https://go.dev/tour/moretypes/15
	m.lines = append(m.lines, line)

	// Rebuild the viewport content from all lines.
	// strings.Join concatenates slice elements with a separator.
	// See: https://pkg.go.dev/strings#Join
	content := strings.Join(m.lines, "\n")
	m.viewport.SetContent(content)

	// Scroll to the bottom so the newest output is visible.
	// GotoBottom() moves the viewport scroll position to the end.
	m.viewport.GotoBottom()
}

// Clear removes all output lines and resets the viewport.
// Call this when starting a fresh installation.
func (m *OutputModel) Clear() {
	// Reset the lines slice. make() creates a new, empty slice.
	// See: https://go.dev/tour/moretypes/13
	m.lines = make([]string, 0)
	m.viewport.SetContent("")
	m.viewport.GotoTop()
}

// SetInstalling toggles the installation progress indicator.
// When installing is true, the progress bar is shown at the top of the tab.
//
// Parameters:
//   - moduleName: The module currently being installed (shown as a label).
//   - installing: Whether installation is in progress.
func (m *OutputModel) SetInstalling(moduleName string, installing bool) {
	m.currentModule = moduleName
	m.isInstalling = installing
}

// SetSize updates the rendering dimensions. Called on terminal resize.
// We also update the viewport and progress bar dimensions.
func (m *OutputModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Update viewport dimensions using setter methods. In Bubbles v2, fields
	// are not exported directly — you use SetWidth/SetHeight methods instead.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport#Model.SetWidth
	vpHeight := height
	if m.isInstalling {
		vpHeight = height - 2
	}
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(vpHeight)

	// Update progress bar width via setter method.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#Model.SetWidth
	m.progress.SetWidth(width)
}

// Update handles messages for the Output tab.
//
// The viewport component handles its own scroll input (page up, page down,
// mouse wheel, arrow keys). We forward all messages to it unconditionally.
// The progress bar also needs to receive its tick messages for animation.
//
// # tea.Batch
//
// When multiple sub-components each return a tea.Cmd, we combine them with
// tea.Batch() so the runtime executes all of them concurrently.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Batch
func (m OutputModel) Update(msg tea.Msg) (OutputModel, tea.Cmd) {
	var cmds []tea.Cmd

	// Forward messages to the viewport. It handles scrolling internally.
	// The viewport's Update() returns a new viewport.Model and an optional
	// command (e.g., for smooth scrolling animations).
	// See: https://pkg.go.dev/charm.land/bubbles/v2/viewport#Model.Update
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	if vpCmd != nil {
		cmds = append(cmds, vpCmd)
	}

	// Forward messages to the progress bar if an install is in progress.
	// In Bubbles v2, progress.Update() returns (progress.Model, tea.Cmd)
	// directly — no type assertion needed, unlike v1 which returned tea.Model.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#Model.Update
	if m.isInstalling {
		var progCmd tea.Cmd
		m.progress, progCmd = m.progress.Update(msg)
		if progCmd != nil {
			cmds = append(cmds, progCmd)
		}
	}

	// Combine all commands. tea.Batch runs them concurrently.
	return m, tea.Batch(cmds...)
}

// View renders the Output tab content.
//
// Layout when installing:
//
//	┌─────────────────────────────────────┐
//	│ Installing: nvim  ████████░░░ 75%   │
//	│─────────────────────────────────────│
//	│ $ stow -v nvim                      │
//	│ LINK: .config/nvim => dotFiles/nvim │
//	│ $ npm install -g neovim             │
//	│ added 1 package...                  │
//	│ ...                 (scrollable)    │
//	└─────────────────────────────────────┘
//
// Layout when idle:
//
//	┌─────────────────────────────────────┐
//	│                                     │
//	│  No installation output yet.        │
//	│  Select a module and press 'i'      │
//	│  to install.                        │
//	│                                     │
//	└─────────────────────────────────────┘
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinVertical
func (m OutputModel) View() string {
	// If there's no output and no install in progress, show a placeholder.
	if len(m.lines) == 0 && !m.isInstalling {
		placeholder := theme.DimText.Render(
			"No installation output yet.\n" +
				"Select a module and press " +
				theme.KeyStyle.Render("i") +
				theme.DimText.Render(" to install."),
		)

		// Centre the placeholder in the available space.
		// lipgloss.Place positions a string within a width×height box.
		// See: https://pkg.go.dev/charm.land/lipgloss/v2#Place
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			placeholder,
		)
	}

	// Build the view from top to bottom.
	var sections []string

	// Show the progress bar header when installing.
	if m.isInstalling {
		// Render the module name label.
		label := theme.Subtitle.Render(
			fmt.Sprintf("Installing: %s", m.currentModule),
		)

		// Render the progress bar. ViewAs(percent) renders the bar at
		// a given percentage without changing the model's internal state.
		// See: https://pkg.go.dev/charm.land/bubbles/v2/progress#Model.ViewAs
		bar := m.progress.View()

		sections = append(sections, label+"\n"+bar)
	}

	// Add the scrollable viewport with all log lines.
	sections = append(sections, m.viewport.View())

	// Stack sections vertically.
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
