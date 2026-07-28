// Package sidebar implements the module list sidebar component for the TUI.
//
// This file defines the Bubble Tea sub-model for the sidebar, which displays
// a searchable, scrollable list of dotfile modules. It follows The Elm Architecture
// (TEA) pattern used by Bubble Tea: Model → Init → Update → View.
//
// # Sub-Model Pattern
//
// In Bubble Tea, complex UIs are built from composable sub-models. Each sub-model
// is a struct that implements its own Init/Update/View cycle. The parent model
// (in our case, the root app model) owns the sub-model and delegates messages
// to it. This is Go's composition-over-inheritance approach.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2
// See: https://go.dev/doc/effective_go#embedding
//
// # Message Passing
//
// Sub-models communicate with their parent by returning tea.Cmd functions that
// produce messages. The parent's Update method receives these messages and can
// react accordingly. This keeps components decoupled.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
package sidebar

import (
	// fmt provides formatted I/O functions like Sprintf (similar to printf in C).
	// See: https://pkg.go.dev/fmt
	"fmt"

	// strings provides functions for manipulating UTF-8 encoded strings.
	// See: https://pkg.go.dev/strings
	"strings"

	// tea is the Bubble Tea framework — our TUI runtime.
	// We alias it as "tea" for brevity (Go allows import aliases).
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// textinput is a pre-built Bubble Tea component for single-line text input.
	// We use it for the search bar.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput
	"charm.land/bubbles/v2/textinput"

	// lipgloss provides CSS-like terminal styling. We use it for layout
	// calculations (measuring rendered string heights).
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	// theme contains our application's shared colour palette and styles.
	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------
// In Bubble Tea, messages (Msg) are plain Go values that represent events.
// Any type can be a message — the convention is to use simple structs.
// The parent model receives these messages and can react to sidebar events.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg

// ModuleSelectedMsg is sent when the user presses Enter on a module.
// The parent model can use this to show the detail panel for the selected module.
type ModuleSelectedMsg struct {
	Name string
}

// CursorChangedMsg is sent whenever the cursor moves to a different module.
// The parent can use this to preview module details as the user navigates.
type CursorChangedMsg struct {
	Name string
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// Model is the Bubble Tea sub-model for the sidebar module list.
//
// In Go, struct fields that start with a lowercase letter are unexported (private
// to the package). This encapsulates the sidebar's internal state and prevents
// other packages from directly mutating it. Instead, they interact through the
// public NewModel(), Update(), and View() methods.
//
// See: https://go.dev/doc/effective_go#names
// See: https://go.dev/ref/spec#Exported_identifiers
type Model struct {
	// items holds the complete unfiltered list of available modules.
	items []ModuleItem

	// filtered holds the current search-filtered list. When the search query
	// is empty, filtered == items. The cursor index refers to this slice.
	filtered []ModuleItem

	// categoryFilter restricts the list to one Category (empty = All).
	categoryFilter string

	// cursor is the zero-based index of the currently highlighted item in
	// the filtered list. It's clamped to [0, len(filtered)-1].
	cursor int

	// width and height define the available space for rendering (in terminal cells).
	// These are set by the parent model, typically from tea.WindowSizeMsg.
	width  int
	height int

	// yOffset is the sidebar's vertical position in the terminal (in rows from the top).
	// Set by the parent model so mouse click Y coordinates can be translated to
	// sidebar-relative positions.
	yOffset int

	// searchMode indicates whether the search input is currently active.
	// When true, key presses are forwarded to the textinput model.
	searchMode bool

	// searchInput is the Bubbles textinput component used for searching.
	// It handles cursor movement, text editing, and rendering of the input field.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput
	searchInput textinput.Model

	// selected is the Name of the currently highlighted module. It's updated
	// whenever the cursor moves, so the parent can read it at any time.
	selected string

	// focused is true when the sidebar panel has keyboard focus. Affects how
	// strongly the cursor item is highlighted vs dimmed siblings.
	focused bool
}

// NewModel creates and returns an initialised sidebar model.
//
// In Go, there are no constructors. By convention, we create "New" functions
// that return a fully initialised struct. This ensures all fields have sensible
// defaults, since Go zero-values (0, "", false, nil) might not always be correct.
//
// Parameters:
//   - items: The full list of modules to display.
//   - width: Available horizontal space in terminal cells.
//   - height: Available vertical space in terminal cells.
//
// See: https://go.dev/doc/effective_go#composite_literals
func NewModel(items []ModuleItem, width int, height int) Model {
	// Initialise the textinput component for the search bar.
	// textinput.New() returns a Model with sensible defaults.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#New
	ti := textinput.New()

	// Placeholder is the greyed-out hint text shown when the input is empty.
	ti.Placeholder = "Type to filter modules..."

	// CharLimit restricts how many characters the user can type.
	// A limit of 50 is generous for module name searches.
	ti.CharLimit = 50

	// Determine the initially selected module name.
	// In Go, we must guard against empty slices to avoid index-out-of-bounds panics.
	// See: https://go.dev/ref/spec#Index_expressions
	initialSelected := ""
	if len(items) > 0 {
		initialSelected = items[0].Name
	}

	// Return the fully initialised model using a struct literal.
	// Go struct literals let you set fields by name — unmentioned fields get
	// their zero values (0, "", false, nil).
	// See: https://go.dev/tour/moretypes/5
	return Model{
		items:       items,
		filtered:    items,
		cursor:      0,
		width:       width,
		height:      height,
		searchMode:  false,
		searchInput: ti,
		selected:    initialSelected,
		focused:     true,
	}
}

// Init is called when the sub-model is first created. It returns an optional
// initial command (tea.Cmd). Since the sidebar has no async initialisation
// (no network calls, no file reads), we return nil.
//
// tea.Cmd is a function type: func() tea.Msg. Returning nil means "do nothing".
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and returns the updated model plus any commands.
//
// This is the core of the Elm Architecture: the model is immutable-ish (we return
// a new copy), and all side effects are expressed as commands (tea.Cmd).
//
// In Bubble Tea v2, the model is passed by value and returned by value. Go copies
// the struct on each call, so modifications inside Update don't affect the caller's
// copy until the return value is used.
//
// The type switch (msg.(type)) is Go's mechanism for inspecting the concrete type
// of an interface value. Each case branch receives a variable of that concrete type.
// See: https://go.dev/tour/methods/16
// See: https://go.dev/doc/effective_go#type_switch
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// We may collect multiple commands to batch together.
	// tea.Batch() combines multiple Cmds into one that runs them concurrently.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#Batch
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// tea.KeyPressMsg is sent when the user presses a key.
	// msg.String() returns a human-readable representation like "j", "enter", "esc".
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#KeyPressMsg
	case tea.KeyPressMsg:
		// When search mode is active, most keys are forwarded to the text input.
		// Only Escape exits search mode.
		if m.searchMode {
			return m.handleSearchModeKey(msg, cmds)
		}

		// Normal (non-search) mode key handling.
		return m.handleNormalModeKey(msg, cmds)

	// tea.MouseClickMsg is sent when a mouse button is clicked.
	// We use it to let users click on sidebar items to select them.
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft && len(m.filtered) > 0 {
			return m.handleMouseClick(msg, cmds)
		}

	// tea.MouseWheelMsg is sent when the scroll wheel is used.
	// We map wheel up/down to cursor up/down for scrolling the list.
	case tea.MouseWheelMsg:
		if msg.Button == tea.MouseWheelUp {
			m = m.moveCursorUp()
			cmds = append(cmds, m.cursorChangedCmd())
		} else if msg.Button == tea.MouseWheelDown {
			m = m.moveCursorDown()
			cmds = append(cmds, m.cursorChangedCmd())
		}
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

// handleSearchModeKey processes key presses when the search input is focused.
//
// This is a helper method extracted from Update() to keep the main switch readable.
// In Go, methods are functions with a special receiver argument that appears before
// the function name. The receiver binds the method to its type.
// See: https://go.dev/tour/methods/1
func (m Model) handleSearchModeKey(msg tea.KeyPressMsg, cmds []tea.Cmd) (Model, tea.Cmd) {
	switch msg.String() {

	// Escape exits the search input but keeps the current filter so the user
	// can j/k through the filtered list. Clear the filter with Esc again
	// in normal mode (see handleNormalModeKey).
	case "esc":
		m.searchMode = false
		m.searchInput.Blur()
		return m, nil

	// Enter in search mode selects the currently highlighted module.
	case "enter":
		if len(m.filtered) > 0 {
			m.searchMode = false
			m.searchInput.Blur()
			m.selected = m.filtered[m.cursor].Name

			// Return a command that produces a ModuleSelectedMsg.
			// tea.Cmd is a function that returns a tea.Msg. The parent's Update
			// will receive this message on the next cycle.
			// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
			return m, func() tea.Msg {
				return ModuleSelectedMsg{Name: m.selected}
			}
		}
		return m, nil

	// Allow cursor navigation even while searching (arrow keys / emacs bindings).
	// Do not bind j/k here — those are typed into the filter query.
	case "up", "ctrl+p":
		m = m.moveCursorUp()
		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)

	case "down", "ctrl+n":
		m = m.moveCursorDown()
		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)

	// All other keys are forwarded to the text input for editing.
	default:
		var tiCmd tea.Cmd

		// Forward the message to the textinput model. It returns an updated
		// model and optionally a command (e.g., for cursor blinking).
		// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Update
		m.searchInput, tiCmd = m.searchInput.Update(msg)
		if tiCmd != nil {
			cmds = append(cmds, tiCmd)
		}

		// After the text changes, re-filter the module list.
		// Value() returns the current text in the input.
		// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Value
		m.reapplyFilters()

		// Reset cursor to the beginning of the new filtered list.
		m.cursor = 0
		if len(m.filtered) > 0 {
			m.selected = m.filtered[0].Name
		} else {
			m.selected = ""
		}

		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)
	}
}

// handleNormalModeKey processes key presses when search mode is NOT active.
func (m Model) handleNormalModeKey(msg tea.KeyPressMsg, cmds []tea.Cmd) (Model, tea.Cmd) {
	switch msg.String() {

	// Esc clears an active text filter (category filter is unchanged).
	case "esc":
		if m.searchInput.Value() == "" {
			return m, tea.Batch(cmds...)
		}
		m.searchInput.SetValue("")
		m.reapplyFilters()
		m.cursor = 0
		if len(m.filtered) > 0 {
			m.selected = m.filtered[0].Name
		} else {
			m.selected = ""
		}
		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)

	// j or down arrow moves the cursor down.
	case "j", "down":
		m = m.moveCursorDown()
		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)

	// k or up arrow moves the cursor up.
	case "k", "up":
		m = m.moveCursorUp()
		cmds = append(cmds, m.cursorChangedCmd())
		return m, tea.Batch(cmds...)

	// Enter selects the current module.
	case "enter":
		if len(m.filtered) > 0 {
			m.selected = m.filtered[m.cursor].Name
			return m, func() tea.Msg {
				return ModuleSelectedMsg{Name: m.selected}
			}
		}
		return m, nil

	// s or / enters search mode and focuses the text input.
	case "s", "/":
		m.searchMode = true

		// Focus() activates the text input so it captures keystrokes.
		// It returns a tea.Cmd (e.g., to start the cursor blink timer).
		// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Focus
		focusCmd := m.searchInput.Focus()
		cmds = append(cmds, focusCmd)
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

// handleMouseClick translates a mouse click into a cursor selection.
// It calculates which item was clicked based on the Y position relative to
// the sidebar's render area, accounting for the search bar and visible range.
func (m Model) handleMouseClick(msg tea.MouseClickMsg, cmds []tea.Cmd) (Model, tea.Cmd) {
	// Translate absolute terminal Y to sidebar-relative Y.
	relY := msg.Y - m.yOffset

	// The search bar (with border) plus a trailing newline occupies
	// several lines at the top of the sidebar.
	sbLines := m.searchBarLines()
	itemY := relY - sbLines
	if itemY < 0 {
		return m, nil
	}

	linesPerItem := 4
	vc := m.visibleItemCount()
	start, _ := m.visibleRange(vc)

	clickedIndex := start + itemY/linesPerItem
	if clickedIndex < 0 || clickedIndex >= len(m.filtered) {
		return m, nil
	}

	prevCursor := m.cursor
	m.cursor = clickedIndex
	m.selected = m.filtered[m.cursor].Name
	cmds = append(cmds, m.cursorChangedCmd())

	// If the clicked item was already selected, treat it as a selection
	// (equivalent to pressing Enter).
	if clickedIndex == prevCursor {
		return m, func() tea.Msg {
			return ModuleSelectedMsg{Name: m.selected}
		}
	}

	return m, tea.Batch(cmds...)
}

// ---------------------------------------------------------------------------
// Cursor Helpers
// ---------------------------------------------------------------------------

// moveCursorUp decrements the cursor, wrapping to the bottom if at the top.
//
// This demonstrates Go's simple if/else control flow. There's no ternary operator
// in Go — explicit if/else is the idiomatic way.
// See: https://go.dev/tour/flowcontrol/5
func (m Model) moveCursorUp() Model {
	if m.cursor > 0 {
		m.cursor--
	} else if len(m.filtered) > 0 {
		// Wrap around to the last item.
		m.cursor = len(m.filtered) - 1
	}

	// Update the selected name to match the new cursor position.
	if len(m.filtered) > 0 {
		m.selected = m.filtered[m.cursor].Name
	}
	return m
}

// moveCursorDown increments the cursor, wrapping to the top if at the bottom.
func (m Model) moveCursorDown() Model {
	if len(m.filtered) == 0 {
		return m
	}

	if m.cursor < len(m.filtered)-1 {
		m.cursor++
	} else {
		// Wrap around to the first item.
		m.cursor = 0
	}

	m.selected = m.filtered[m.cursor].Name
	return m
}

// cursorChangedCmd returns a tea.Cmd that produces a CursorChangedMsg.
// This notifies the parent model that the user is hovering over a different module.
//
// In Go, functions are first-class values. A tea.Cmd is simply a func() tea.Msg.
// We use a closure (anonymous function) that captures m.selected from the
// enclosing scope.
// See: https://go.dev/tour/moretypes/25
func (m Model) cursorChangedCmd() tea.Cmd {
	// Capture the selected name in a local variable to avoid closure issues.
	// If we captured m.selected directly, the closure would reference the Model's
	// field, which might change by the time the Cmd executes.
	name := m.selected
	return func() tea.Msg {
		return CursorChangedMsg{Name: name}
	}
}

// ---------------------------------------------------------------------------
// Public Accessors
// ---------------------------------------------------------------------------
// These methods allow the parent model to read sidebar state without exposing
// internal fields. This is Go's approach to encapsulation.
// See: https://go.dev/doc/effective_go#Getters

// Selected returns the name of the currently highlighted module.
func (m Model) Selected() string {
	return m.selected
}

// IsSearching reports whether the sidebar's search input is currently active.
// The parent model uses this to avoid intercepting keys meant for search.
func (m Model) IsSearching() bool {
	return m.searchMode
}

// SetSize updates the sidebar's available dimensions. The parent calls this
// when the terminal is resized (tea.WindowSizeMsg).
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetYOffset stores the sidebar's vertical screen position so mouse clicks
// can be translated from absolute terminal coordinates to sidebar-relative ones.
func (m *Model) SetYOffset(y int) {
	m.yOffset = y
}

// ActivateSearch enables search mode and focuses the text input.
// This allows the parent to trigger search from any focus context.
func (m *Model) ActivateSearch() tea.Cmd {
	m.searchMode = true
	m.focused = true
	return m.searchInput.Focus()
}

// SetFocused marks whether the sidebar panel has keyboard focus.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// HasFilter reports whether a text search filter is active.
func (m Model) HasFilter() bool {
	return m.searchInput.Value() != ""
}

// SetInstalled updates the install status of a module by name.
// This is called after install/uninstall operations so the sidebar
// reflects the current state without requiring a restart.
func (m *Model) SetInstalled(name string, installed bool) {
	for i := range m.items {
		if m.items[i].Name == name {
			m.items[i].Installed = installed
			break
		}
	}
	for i := range m.filtered {
		if m.filtered[i].Name == name {
			m.filtered[i].Installed = installed
			break
		}
	}
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

// View renders the sidebar as a styled string.
//
// In Bubble Tea, View() must be a pure function — it reads from the model but
// never modifies it or performs I/O. The returned string contains ANSI escape
// codes for colours and styling, which the terminal interprets.
//
// View returns a plain string (not tea.View) because this is a sub-model.
// Only the root model returns tea.View.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Model
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.renderSearchBar())
	b.WriteString("\n")

	if len(m.filtered) == 0 {
		noResults := theme.DimText.Render("  No modules match your search.")
		b.WriteString(noResults)
		b.WriteString("\n")
	} else {
		b.WriteString(m.renderModuleList())
	}

	// Hard-clip to the size the parent allocated so panel borders stay intact.
	return theme.Clip(b.String(), m.width, m.height)
}

// renderSearchBar renders the search input area at the top of the sidebar.
func (m Model) renderSearchBar() string {
	catLabel := "All"
	if m.categoryFilter != "" {
		catLabel = m.categoryFilter
	}
	catLine := theme.DimText.Render(fmt.Sprintf("  Category: %s  (c to change)", catLabel))

	if m.searchMode {
		return catLine + "\n" + m.searchInput.View()
	}

	hint := theme.DimText.Render("  / search")
	if m.searchInput.Value() != "" {
		hint = theme.NormalText.Render("  filter: "+m.searchInput.Value()) +
			theme.DimText.Render("  (esc clear · / edit)")
	}
	return catLine + "\n" + hint
}

// CategoryFilter returns the active category filter (empty = All).
func (m Model) CategoryFilter() string {
	return m.categoryFilter
}

// SetItems replaces the module list and re-applies active filters.
func (m *Model) SetItems(items []ModuleItem) {
	m.items = items
	m.reapplyFilters()
	m.cursor = 0
	if len(m.filtered) > 0 {
		m.selected = m.filtered[0].Name
	} else {
		m.selected = ""
	}
}

// SetCategoryFilter sets the category filter and refreshes the list.
// Pass "" for All categories.
func (m *Model) SetCategoryFilter(category string) {
	m.categoryFilter = category
	m.reapplyFilters()
	m.cursor = 0
	if len(m.filtered) > 0 {
		m.selected = m.filtered[0].Name
	} else {
		m.selected = ""
	}
}

// reapplyFilters applies category then search filters to items → filtered.
func (m *Model) reapplyFilters() {
	list := FilterByCategory(m.items, m.categoryFilter)
	m.filtered = FilterModules(list, m.searchInput.Value())
}

// renderModuleList renders the scrollable list of module items.
//
// This method handles the "viewport" logic: if the list is taller than the
// available height, it shows a window of items centred around the cursor.
func (m Model) renderModuleList() string {
	var b strings.Builder

	visibleCount := m.visibleItemCount()
	start, end := m.visibleRange(visibleCount)

	for i := start; i < end; i++ {
		item := m.filtered[i]
		isActive := i == m.cursor

		b.WriteString(m.renderItem(item, isActive))

		if i < end-1 {
			b.WriteString("\n")
		}
	}

	if len(m.filtered) > visibleCount {
		scrollInfo := theme.DimText.Render(
			fmt.Sprintf("  ↕ %d/%d", m.cursor+1, len(m.filtered)),
		)
		b.WriteString("\n")
		b.WriteString(scrollInfo)
	}

	return b.String()
}

// visibleRange calculates the start (inclusive) and end (exclusive) indices
// for the window of items to display, centred around the cursor.
//
// This function returns two int values — Go supports multiple return values,
// which is idiomatic for returning related results without defining a struct.
// See: https://go.dev/doc/effective_go#multiple-returns
func (m Model) visibleRange(visibleCount int) (int, int) {
	total := len(m.filtered)

	// If all items fit, show everything.
	if visibleCount >= total {
		return 0, total
	}

	// Centre the window around the cursor.
	// The integer division by 2 gives us half the window size.
	half := visibleCount / 2
	start := m.cursor - half

	// Clamp start to valid bounds [0, total - visibleCount].
	if start < 0 {
		start = 0
	}
	if start > total-visibleCount {
		start = total - visibleCount
	}

	end := start + visibleCount
	return start, end
}

// renderItem renders a single module item as a bordered box.
//
// The box contains:
//   - Line 1: Icon + Name (bold)
//   - Line 2: Description (dimmed)
//   - Line 3: Install status indicator (✓ green / ✗ red)
func (m Model) renderItem(item ModuleItem, isActive bool) string {
	icon := item.Icon
	if icon == "" {
		icon = theme.GetModuleIcon(item.Name)
	}

	desc := item.Description
	maxDescLen := m.contentWidth() - 4
	if maxDescLen > 0 && len(desc) > maxDescLen {
		desc = desc[:maxDescLen-1] + "…"
	}

	var nameLine, descLine, statusLine string
	if isActive {
		// Nested styles must set the same Background — otherwise ANSI resets
		// from Foreground-only spans punch holes in the highlight fill.
		bg := lipgloss.NewStyle().
			Background(theme.ColorSurface).
			ColorWhitespace(true)
		nameLine = bg.Foreground(theme.ColorPink).Bold(true).
			Render(fmt.Sprintf("▸ %s  %s", icon, item.Name))
		descLine = bg.Foreground(theme.ColorForeground).Render(desc)
		if item.Installed {
			statusLine = bg.Foreground(theme.ColorGreen).
				Render(theme.IconInstalled + " installed")
		} else {
			statusLine = bg.Foreground(theme.ColorForegroundDim).
				Render(theme.IconNotInstalled + " not installed")
		}
	} else {
		nameLine = lipgloss.NewStyle().
			Foreground(theme.ColorForegroundDim).
			Faint(true).
			Render(fmt.Sprintf("  %s  %s", icon, item.Name))
		descLine = lipgloss.NewStyle().
			Foreground(theme.ColorForegroundMuted).
			Faint(true).
			Render(desc)
		if item.Installed {
			statusLine = theme.StatusInstalled.Render() + " " + theme.SuccessText.Render("installed")
		} else {
			statusLine = lipgloss.NewStyle().Faint(true).Render(
				theme.StatusNotInstalled.Render() + " " + theme.DimText.Render("not installed"),
			)
		}
	}

	content := strings.Join([]string{nameLine, descLine, statusLine}, "\n")

	var style lipgloss.Style
	if isActive {
		style = theme.SidebarItemActive
		if !m.focused {
			style = theme.SidebarItemActive.BorderForeground(theme.ColorPurple)
		}
	} else if item.Installed {
		style = theme.SidebarItemInstalled
	} else {
		style = theme.SidebarItem
	}

	return style.
		Width(m.contentWidth()).
		Render(content)
}

// contentWidth calculates the usable width inside the sidebar, accounting
// for the border characters (1 char on each side = 2 total) and padding.
//
// This is a private helper method (lowercase name). In Go, unexported methods
// are only accessible within the same package.
// See: https://go.dev/doc/effective_go#names
func (m Model) contentWidth() int {
	// Subtract 4 for border (2) + padding (2) from SidebarItem style.
	w := m.width - 4
	if w < 10 {
		w = 10
	}
	return w
}

// searchBarLines returns the total number of terminal lines occupied by the
// search bar area, including the trailing newline added in View().
func (m Model) searchBarLines() int {
	return lipgloss.Height(m.renderSearchBar()) + 1
}

// visibleItemCount returns how many module items fit in the visible area.
//
// Each item is a rounded box: top border + 3 content lines + bottom border = 5,
// plus a 1-line gap between items. The search bar height is measured, and one
// line is reserved for the scroll indicator when the list is longer than the window.
func (m Model) visibleItemCount() int {
	const (
		itemLines = 5 // bordered module card
		itemGap   = 1 // newline between cards
	)

	avail := m.height - m.searchBarLines()
	if avail < itemLines {
		return 1
	}

	// Reserve a footer line when scrolling will be needed. We don't know the
	// final count yet, so reserve if there is more than one module overall.
	if len(m.filtered) > 1 {
		avail--
	}
	if avail < itemLines {
		return 1
	}

	// n*itemLines + (n-1)*itemGap <= avail  →  n <= (avail+itemGap)/(itemLines+itemGap)
	n := (avail + itemGap) / (itemLines + itemGap)
	if n < 1 {
		n = 1
	}
	if n > len(m.filtered) {
		n = len(m.filtered)
	}
	return n
}
