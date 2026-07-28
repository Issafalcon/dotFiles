// Package popup provides reusable modal overlay components for the TUI.
//
// In terminal UIs, popups (or modals/dialogs) are "floating" elements that
// appear on top of the main content — similar to modal dialogs in web apps.
// This package implements several popup types:
//
//   - Model: A generic popup overlay (the base for all others)
//   - ConfirmModel: A Yes/No confirmation dialog
//   - InputModel: A text input dialog
//   - HelpModel: A keyboard shortcuts overlay
//
// Each popup embeds the base Model using Go's struct embedding feature, which
// provides a form of composition (not inheritance — Go has no inheritance).
// See: https://go.dev/doc/effective_go#embedding
//
// The popups are designed to be used as sub-models within a Bubble Tea
// application. They are NOT top-level tea.Model implementations — they use
// concrete return types instead of interfaces, which is the standard pattern
// for sub-components in Bubble Tea.
// See: https://pkg.go.dev/charm.land/bubbletea/v2
package popup

import (
	"strings"

	// lipgloss is a CSS-like terminal styling library. We alias the import
	// to "lipgloss" for readability — the full module path is the v2 version.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Model — Generic Popup Overlay
// ---------------------------------------------------------------------------

// Model is the base struct for all popup overlays. It stores the popup's
// content, dimensions, visibility state, and title.
//
// In Go, struct fields that start with a lowercase letter are "unexported"
// (private to the package). This enforces encapsulation — external packages
// must use the exported methods (Show, Hide, etc.) to interact with the model.
// See: https://go.dev/doc/effective_go#names
type Model struct {
	// title is displayed at the top of the popup, styled with theme.Title.
	title string

	// content is the body text rendered inside the popup border.
	content string

	// width and height define the inner dimensions of the popup box
	// (excluding border and padding, which are added by theme.PopupStyle).
	width  int
	height int

	// visible controls whether the popup is shown. When false, Render()
	// returns an empty string and the popup is effectively hidden.
	visible bool
}

// NewPopup creates a new popup Model with the given title, content, and
// dimensions. The popup starts hidden — call Show() to make it visible.
//
// In Go, constructor functions are conventionally named New<Type> or
// New<Package>. They return a value (not a pointer) for small structs,
// which is efficient because Go copies small structs cheaply on the stack.
// See: https://go.dev/doc/effective_go#composite_literals
func NewPopup(title string, content string, width, height int) Model {
	return Model{
		title:   title,
		content: content,
		width:   width,
		height:  height,
		visible: false,
	}
}

// Show returns a copy of the Model with visible set to true.
//
// This uses a value receiver (m Model) rather than a pointer receiver
// (*Model), so the method operates on a copy of the struct. This is
// intentional — it follows the immutable update pattern used throughout
// Bubble Tea, where state changes produce new values instead of mutating
// existing ones.
// See: https://go.dev/tour/methods/4
// See: https://go.dev/doc/effective_go#methods
func (m Model) Show() Model {
	m.visible = true
	return m
}

// Hide returns a copy of the Model with visible set to false.
// Same immutable-update pattern as Show().
func (m Model) Hide() Model {
	m.visible = false
	return m
}

// IsVisible reports whether the popup is currently shown.
//
// By Go convention, boolean getters are named without a "Get" prefix.
// Predicate methods often start with "Is", "Has", or "Can".
// See: https://go.dev/doc/effective_go#Getters
func (m Model) IsVisible() bool {
	return m.visible
}

// SetContent returns a copy of the Model with updated content.
// This allows parent components to update what the popup displays.
func (m Model) SetContent(content string) Model {
	m.content = content
	return m
}

// Render produces the final string output for the popup, centered on screen.
// If the popup is not visible, it returns an empty string.
//
// Parameters:
//   - screenWidth:  the full terminal width in columns
//   - screenHeight: the full terminal height in rows
//
// The rendering pipeline is:
//  1. Build the popup body (title + separator + content)
//  2. Apply theme.PopupStyle (double border, padding, background)
//  3. Center the styled box on the screen using lipgloss.Place
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2#Place
func (m Model) Render(screenWidth, screenHeight int) string {
	if !m.visible {
		return ""
	}
	return renderPopup(m.title, m.content, m.width, m.height, screenWidth, screenHeight)
}

// ---------------------------------------------------------------------------
// Internal Helpers
// ---------------------------------------------------------------------------

// renderPopup is a package-level helper shared by all popup types. It wraps
// the given content in a themed popup box and centers it on the screen.
//
// This is an unexported function (lowercase first letter), meaning it can
// only be called from within the popup package. Sub-model types (ConfirmModel,
// InputModel, HelpModel) call this to avoid duplicating rendering logic.
// See: https://go.dev/doc/effective_go#names
func renderPopup(title, content string, width, height, screenWidth, screenHeight int) string {
	// --- Step 1: Build the popup body ---
	titleStr := theme.Title.Render(title)

	// Create a horizontal rule to visually separate the title from content.
	// strings.Repeat repeats a string N times — here we create a line of
	// box-drawing characters. We subtract 4 to account for PopupStyle's
	// horizontal padding (2 on each side).
	// See: https://pkg.go.dev/strings#Repeat
	separatorWidth := width - 4
	if separatorWidth < 1 {
		separatorWidth = 1
	}
	separator := theme.DimText.Render(strings.Repeat("─", separatorWidth))

	// Compose the full body: title, separator, blank line, then content.
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStr,
		separator,
		"",
		content,
	)

	// --- Step 2: Apply popup styling ---
	// theme.PopupStyle applies a double border (╔═╗║ ║╚═╝), pink border
	// colour, internal padding, and a dark background. Width() and Height()
	// set the inner content area dimensions.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#Style.Width
	popup := theme.PopupStyle.
		Width(width).
		Height(height).
		Render(body)

	// --- Step 3: Center on screen ---
	// lipgloss.Place positions a string within a larger area. Here we centre
	// the popup both horizontally and vertically within the full terminal.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#Place
	return lipgloss.Place(
		screenWidth, screenHeight,
		lipgloss.Center, lipgloss.Center,
		popup,
	)
}
