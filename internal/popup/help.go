// This file implements the help overlay popup that displays keyboard shortcuts.
//
// The help popup is a read-only overlay — it doesn't accept input beyond
// a dismiss key. It shows all available keyboard shortcuts grouped by
// category (Navigation, Actions, UI, App) in a two-column layout.
//
// # Two-Column Layout with Lip Gloss
//
// Terminal UIs don't have CSS Grid or Flexbox, so we build column layouts
// manually. Each row is a pair of styled strings (key + description) joined
// horizontally with lipgloss.JoinHorizontal. Rows are then stacked vertically.
// Fixed-width rendering ensures columns stay aligned.
// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinHorizontal
package popup

import (
	"strings"

	// lipgloss provides terminal styling and layout utilities.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// HelpBinding represents a single keyboard shortcut entry.
// It pairs a key (e.g., "↑/k") with its description (e.g., "Move up").
//
// In Go, exported struct fields (capitalised) can be accessed from other
// packages. This lets the parent application provide custom key bindings.
// See: https://go.dev/ref/spec#Exported_identifiers
type HelpBinding struct {
	Key         string
	Description string
}

// helpSection groups related key bindings under a category heading.
// This is unexported (lowercase) because it's only used internally
// to organise the rendering of the help overlay.
type helpSection struct {
	title    string
	bindings []HelpBinding
}

// HelpModel is the help overlay popup. It embeds the base popup Model
// and stores the list of key bindings to display.
type HelpModel struct {
	// Embedded base popup — provides title, dimensions, visibility, Render().
	// See: https://go.dev/doc/effective_go#embedding
	Model

	// keyBindings is the flat list of all shortcuts. It's stored so that
	// callers can provide custom bindings via NewHelpPopup().
	keyBindings []HelpBinding
}

// ---------------------------------------------------------------------------
// Constructor
// ---------------------------------------------------------------------------

// NewHelpPopup creates a help overlay showing keyboard shortcuts.
//
// If bindings is nil, a default set of shortcuts is used. This follows the
// "sensible defaults with optional override" pattern common in Go libraries.
//
// The nil check uses Go's zero value: an uninitialized slice is nil, and
// len(nil) == 0. However, a non-nil empty slice ([]HelpBinding{}) is
// different from nil — we only use defaults when explicitly nil.
// See: https://go.dev/doc/effective_go#allocation_new
// See: https://go.dev/tour/moretypes/12
func NewHelpPopup(bindings []HelpBinding) HelpModel {
	if bindings == nil {
		bindings = defaultBindings()
	}

	return HelpModel{
		Model:       NewPopup("⌨ Keyboard Shortcuts", "", 55, 22).Show(),
		keyBindings: bindings,
	}
}

// defaultBindings returns the standard set of keyboard shortcuts for the app.
//
// This is a package-level function (not a method) because it doesn't depend
// on any receiver. In Go, free functions are perfectly normal and preferred
// over unnecessary methods.
// See: https://go.dev/doc/effective_go#functions
func defaultBindings() []HelpBinding {
	return []HelpBinding{
		// Navigation
		{Key: "↑/k", Description: "Move up"},
		{Key: "↓/j", Description: "Move down"},
		{Key: "enter", Description: "Select item"},
		// Actions
		{Key: "i", Description: "Install module"},
		{Key: "d", Description: "Uninstall module"},
		{Key: "o", Description: "Open URL in browser"},
		// UI
		{Key: "tab", Description: "Switch tabs"},
		{Key: "s", Description: "Search modules"},
		{Key: "?", Description: "Toggle help"},
		// App
		{Key: "q", Description: "Quit application"},
		{Key: "esc", Description: "Cancel / close popup"},
	}
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

// View builds the inner content of the help overlay as a formatted string.
//
// The layout is a two-column table grouped by category:
//
//	 Navigation
//	   ↑/k     Move up
//	   ↓/j     Move down
//	   enter   Select item
//
//	 Actions
//	   i       Install module
//	   ...
//
// Each key is rendered with theme.KeyStyle (bold pink) at a fixed width,
// and each description with theme.DescStyle (dim text). The fixed width
// ensures the description column lines up neatly.
func (m HelpModel) View() string {
	// Organise the flat binding list into categorised sections.
	// We use the helpSection type to group bindings by category.
	sections := buildSections(m.keyBindings)

	// strings.Builder efficiently accumulates the output string.
	// See: https://pkg.go.dev/strings#Builder
	var b strings.Builder

	// keyColWidth is the fixed character width for the key column.
	// lipgloss.Width doesn't measure raw runes — it accounts for
	// ANSI escape sequences (colours) and wide Unicode characters.
	// We use a fixed width here for simplicity and consistency.
	const keyColWidth = 12

	// Iterate over each section and render its title + bindings.
	for i, section := range sections {
		// --- Section title ---
		// Rendered with Subtitle style (bold cyan) for visual hierarchy.
		b.WriteString(theme.Subtitle.Render(section.title))
		b.WriteString("\n")

		// --- Key-description rows ---
		for _, binding := range section.bindings {
			// Render the key with a fixed width so descriptions align.
			// lipgloss.NewStyle().Width(n) pads or truncates to exactly
			// n columns, creating consistent column alignment.
			// See: https://pkg.go.dev/charm.land/lipgloss/v2#Style.Width
			keyCell := theme.KeyStyle.
				Width(keyColWidth).
				Render(binding.Key)

			descCell := theme.DescStyle.Render(binding.Description)

			// Join the key and description horizontally on one line.
			row := lipgloss.JoinHorizontal(lipgloss.Top, "  ", keyCell, descCell)
			b.WriteString(row)
			b.WriteString("\n")
		}

		// Add a blank line between sections (but not after the last one).
		// len(sections)-1 gives the index of the last element.
		if i < len(sections)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Render produces the final help overlay: content wrapped in the themed
// popup box, centered on the terminal screen.
func (m HelpModel) Render(screenWidth, screenHeight int) string {
	if !m.Model.IsVisible() {
		return ""
	}
	return renderPopup(
		m.Model.title, m.View(),
		m.Model.width, m.Model.height,
		screenWidth, screenHeight,
	)
}

// ---------------------------------------------------------------------------
// Internal Helpers
// ---------------------------------------------------------------------------

// buildSections groups a flat list of HelpBindings into categorised sections.
//
// The grouping is based on the default layout:
//   - Indices 0–2  → Navigation
//   - Indices 3–5  → Actions
//   - Indices 6–8  → UI
//   - Indices 9–10 → App
//
// If the binding list is shorter than expected (custom bindings), we
// gracefully handle it by slicing only up to the available length.
//
// The min() built-in function (added in Go 1.21) returns the smaller of
// two values. It prevents out-of-bounds panics when the slice is shorter
// than a boundary index.
// See: https://pkg.go.dev/builtin#min
func buildSections(bindings []HelpBinding) []helpSection {
	// safeSlice returns bindings[start:end], clamped to the actual slice length.
	// This is an anonymous function (closure) that captures `bindings`.
	// See: https://go.dev/tour/moretypes/25
	safeSlice := func(start, end int) []HelpBinding {
		// Clamp indices to the slice bounds to avoid a runtime panic.
		if start >= len(bindings) {
			return nil
		}
		if end > len(bindings) {
			end = len(bindings)
		}
		return bindings[start:end]
	}

	// Define the section boundaries. Each section spans a range of indices
	// in the flat bindings list.
	//
	// A slice literal []helpSection{...} creates a new slice with the
	// given elements. This is similar to array literals in other languages.
	// See: https://go.dev/tour/moretypes/9
	sections := []helpSection{
		{title: "Navigation", bindings: safeSlice(0, 3)},
		{title: "Actions", bindings: safeSlice(3, 6)},
		{title: "UI", bindings: safeSlice(6, 9)},
		{title: "App", bindings: safeSlice(9, 11)},
	}

	// Filter out sections that have no bindings (in case the list was short).
	// We build a new slice by appending only non-empty sections.
	//
	// result is declared with var, giving it the zero value (nil slice).
	// append() creates a new backing array when needed and returns the
	// updated slice header.
	// See: https://go.dev/doc/effective_go#slices
	// See: https://pkg.go.dev/builtin#append
	var result []helpSection
	for _, s := range sections {
		if len(s.bindings) > 0 {
			result = append(result, s)
		}
	}

	return result
}
