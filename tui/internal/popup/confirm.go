// This file implements the install confirmation dialog popup.
//
// When a user chooses to install a module, this dialog appears showing what
// will be installed and asking for confirmation. It's a simple Yes/No dialog
// with keyboard navigation between the two buttons.
//
// # Message-Based Communication
//
// In Bubble Tea, components communicate by returning tea.Cmd functions that
// produce messages. When the user confirms or cancels, this component returns
// a command that yields either ConfirmYesMsg or ConfirmNoMsg. The parent
// component handles these messages in its own Update() method.
//
// This pattern decouples the popup from the rest of the app — the popup
// doesn't need to know what happens after confirmation.
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
package popup

import (
	"fmt"
	"strings"

	// tea is the Bubble Tea framework. We alias the import for brevity.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// lipgloss provides CSS-like terminal styling.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------
// In Bubble Tea, messages are simple Go types (usually structs) that carry
// data about events. They flow through the Update() method via the tea.Msg
// interface. Any type satisfies tea.Msg because it's defined as:
//
//	type Msg interface{}
//
// By convention, message type names end with "Msg".
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg

// ConfirmAction indicates whether the confirm dialog is for install or uninstall.
// Using a named string type instead of a raw string gives us type safety — the
// compiler ensures we only use defined constants, not arbitrary strings.
// See: https://go.dev/ref/spec#Type_definitions
type ConfirmAction string

const (
	// ActionInstall means the user is confirming an installation.
	ActionInstall ConfirmAction = "install"
	// ActionUninstall means the user is confirming an uninstallation.
	ActionUninstall ConfirmAction = "uninstall"
)

// ConfirmYesMsg is sent when the user confirms the action (install or uninstall).
// It carries the module name and action so the parent knows what was confirmed.
type ConfirmYesMsg struct {
	ModuleName string
	Action     ConfirmAction
}

// ConfirmNoMsg is sent when the user cancels or declines.
// It has no fields because the parent just needs to know the action was
// cancelled — it doesn't matter which button was active.
type ConfirmNoMsg struct{}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// ConfirmModel is the confirmation dialog. It embeds popup.Model to inherit
// the base popup's fields and methods (title, dimensions, visibility, etc.).
//
// Go struct embedding is NOT inheritance — it's composition. The embedded
// Model's exported methods are "promoted" to ConfirmModel, so you can call
// m.IsVisible() directly. However, if ConfirmModel defines a method with
// the same name, the embedded method is "shadowed" (not overridden).
// See: https://go.dev/doc/effective_go#embedding
// See: https://go.dev/ref/spec#Struct_types
type ConfirmModel struct {
	// Model is the embedded base popup (provides title, dimensions, visibility).
	// Because it's embedded without a field name, its methods are promoted.
	Model

	// moduleName is the name of the module being installed/uninstalled (e.g., "nvim").
	// Stored so we can include it in ConfirmYesMsg.
	moduleName string

	// action distinguishes install vs uninstall confirmations.
	// The parent uses this in ConfirmYesMsg to know which flow to trigger.
	action ConfirmAction

	// items lists what will be installed or uninstalled.
	// In Go, slices ([]string) are dynamic arrays — they grow as needed.
	// See: https://go.dev/doc/effective_go#slices
	items []string

	// subtitle is the descriptive text shown above the items list.
	// Customised per action: "will be installed" vs "will be uninstalled".
	subtitle string

	// confirmed tracks whether the user has made a selection.
	confirmed bool

	// cursor tracks which button is highlighted: 0 = Yes, 1 = No.
	cursor int
}

// NewConfirmDialog creates a new confirmation dialog for installing a module.
//
// Parameters:
//   - moduleName: the display name of the module (used in the title)
//   - deps: a list of items that will be installed
//
// The dialog starts visible (Show() is called on the base popup) with the
// cursor on "Yes" (index 0).
//
// fmt.Sprintf formats a string using printf-style verbs. %s inserts a string.
// See: https://pkg.go.dev/fmt#Sprintf
func NewConfirmDialog(moduleName string, deps []string) ConfirmModel {
	// Build the popup title with the module name interpolated.
	title := fmt.Sprintf("Install %s?", moduleName)

	return ConfirmModel{
		// NewPopup creates the base model; Show() makes it immediately visible.
		// Method chaining works because Show() returns a new Model value.
		Model:      NewPopup(title, "", 50, 15).Show(),
		moduleName: moduleName,
		action:     ActionInstall,
		items:      deps,
		subtitle:   "The following will be installed:",
		confirmed:  false,
		cursor:     0, // Start on "Yes"
	}
}

// NewUninstallDialog creates a new confirmation dialog for uninstalling a module.
//
// Parameters:
//   - moduleName: the display name of the module to uninstall
//   - items: a list describing what will be removed (commands, stow links, etc.)
//
// This mirrors NewConfirmDialog but uses uninstall-specific wording and the
// ActionUninstall action so the parent handler knows to run uninstall commands.
func NewUninstallDialog(moduleName string, items []string) ConfirmModel {
	title := fmt.Sprintf("Uninstall %s?", moduleName)

	return ConfirmModel{
		Model:      NewPopup(title, "", 50, 15).Show(),
		moduleName: moduleName,
		action:     ActionUninstall,
		items:      items,
		subtitle:   "The following will be removed:",
		confirmed:  false,
		cursor:     0, // Start on "Yes"
	}
}

// Update handles keyboard input for the confirmation dialog.
//
// This is the core of the Elm Architecture's "update" step. It receives a
// tea.Msg (which could be any type) and returns the updated model plus an
// optional command (tea.Cmd).
//
// The type switch `switch msg := msg.(type)` is a Go idiom for handling
// interface values. It checks the concrete type of msg and binds it to the
// variable in each case branch. This is called a "type assertion switch".
// See: https://go.dev/tour/methods/16
// See: https://go.dev/doc/effective_go#type_switch
//
// Note: this returns (ConfirmModel, tea.Cmd), NOT (tea.Model, tea.Cmd).
// Sub-components return their concrete type; only the root model needs to
// satisfy the tea.Model interface.
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	// Type-switch on the incoming message to determine what happened.
	switch msg := msg.(type) {

	// tea.KeyPressMsg is sent when the user presses a key.
	// msg.String() returns a human-readable key name like "left", "enter", etc.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#KeyPressMsg
	case tea.KeyPressMsg:
		switch msg.String() {

		// Navigate left to the "Yes" button.
		// We support both arrow keys and Vim-style h/l for accessibility.
		case "left", "h":
			m.cursor = 0

		// Navigate right to the "No" button.
		case "right", "l":
			m.cursor = 1

		// Confirm the current selection.
		case "enter":
			if m.cursor == 0 {
				// User selected "Yes" — mark as confirmed and send a message.
				m.confirmed = true

				// A tea.Cmd is a function that returns a tea.Msg. Bubble Tea
				// runs commands asynchronously and feeds the resulting message
				// back into Update(). Here we use an anonymous function
				// (closure) that captures m.moduleName from the outer scope.
				// See: https://go.dev/tour/moretypes/25
				// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
				return m, func() tea.Msg {
					return ConfirmYesMsg{ModuleName: m.moduleName, Action: m.action}
				}
			}
			// User selected "No" — send a cancel message.
			return m, func() tea.Msg {
				return ConfirmNoMsg{}
			}

		// Escape always cancels, regardless of cursor position.
		case "esc":
			return m, func() tea.Msg {
				return ConfirmNoMsg{}
			}
		}
	}

	// For any unhandled message, return the model unchanged with no command.
	// Returning nil as the command means "do nothing" — no new message will
	// be produced.
	return m, nil
}

// View builds the inner content of the confirmation dialog as a string.
//
// This does NOT include the popup border/centering — that's handled by
// Render(). View() only produces the dialog's body content:
//   - A subtitle ("The following will be installed:")
//   - A bulleted list of items
//   - [Yes] and [No] buttons with the active one highlighted
//
// strings.Builder is an efficient way to build strings incrementally.
// Unlike string concatenation (which creates a new string each time),
// Builder minimises memory allocations by writing to an internal buffer.
// See: https://pkg.go.dev/strings#Builder
func (m ConfirmModel) View() string {
	// strings.Builder accumulates string parts efficiently.
	var b strings.Builder

	// --- Items list ---
	b.WriteString(theme.Subtitle.Render(m.subtitle))
	b.WriteString("\n\n")

	// range iterates over the slice, yielding (index, value) pairs.
	// We use _ to discard the index since we only need the value.
	// See: https://go.dev/tour/moretypes/16
	for _, item := range m.items {
		b.WriteString(theme.NormalText.Render("  • " + item))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// --- Buttons ---
	// Create separate styles for the Yes and No buttons. The active button
	// gets a bright coloured background; the inactive one is dimmed.
	//
	// lipgloss.NewStyle() creates a blank style. Method chaining builds it up.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#NewStyle
	yesStyle := lipgloss.NewStyle().Padding(0, 3)
	noStyle := lipgloss.NewStyle().Padding(0, 3)

	// Apply the appropriate styling based on which button is focused.
	if m.cursor == 0 {
		// Yes is active — green background, dark foreground for contrast.
		yesStyle = yesStyle.
			Bold(true).
			Foreground(theme.ColorBackground).
			Background(theme.ColorGreen)
		noStyle = noStyle.
			Foreground(theme.ColorForegroundDim)
	} else {
		// No is active — red background, dark foreground.
		yesStyle = yesStyle.
			Foreground(theme.ColorForegroundDim)
		noStyle = noStyle.
			Bold(true).
			Foreground(theme.ColorBackground).
			Background(theme.ColorRed)
	}

	// lipgloss.JoinHorizontal places strings side by side.
	// The first argument is the vertical alignment of the joined pieces.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinHorizontal
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		yesStyle.Render(" Yes "),
		"  ", // Spacer between buttons
		noStyle.Render(" No "),
	)

	b.WriteString(buttons)

	return b.String()
}

// Render produces the final string output: the dialog content wrapped in the
// themed popup box, centered on the screen.
//
// This calls View() to build the inner content, then delegates to the
// package-level renderPopup() helper for the border and positioning.
//
// Note: this method "shadows" the embedded Model.Render() method. In Go,
// when an embedding struct defines a method with the same name as the
// embedded type, the outer method takes precedence. The embedded method
// is still accessible via m.Model.Render() if needed.
// See: https://go.dev/ref/spec#Selectors
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
