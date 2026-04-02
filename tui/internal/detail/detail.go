// Package detail implements the right-side detail panel of the DotFiles TUI.
//
// The detail panel uses a tabbed interface with three views:
//
//   - Overview: Module information, dependencies, and their install status
//   - Output: Live scrollable log of installation commands and output
//   - Configuration: Module-specific settings the user can customise
//
// # Tab Container Pattern
//
// This file implements the tab container — the outer shell that renders a tab
// bar at the top and delegates to whichever sub-model is currently active.
// Each tab is a separate Model struct (OverviewModel, OutputModel, ConfigModel)
// with its own Update/View cycle. The container routes messages to whichever
// tab is active.
//
// # Enum Pattern in Go
//
// Go doesn't have a built-in enum keyword. We create enumerations by declaring
// a new named type (type Tab int) and a const block with iota. iota auto-
// increments from 0, giving each constant a unique integer value.
//
// See: https://go.dev/ref/spec#Iota
// See: https://go.dev/ref/spec#Type_definitions
//
// # The Elm Architecture (TEA) for Sub-Models
//
// In Bubble Tea, you can nest models inside each other. The parent model
// owns the child, calls its Update() with messages, and calls its View()
// to get the rendered string. The child doesn't implement tea.Model directly —
// it just has the same method signatures but returns its concrete type
// instead of tea.Model.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2
// See: https://go.dev/doc/effective_go#embedding
package detail

import (
	"strings"

	// tea is the Bubble Tea framework — the core TUI library.
	// We alias the import for brevity, which is idiomatic in the Charm ecosystem.
	// See: https://go.dev/ref/spec#Import_declarations
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// lipgloss provides CSS-like styling for terminal output.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	// theme contains all shared colour and style definitions for the app.
	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Tab Enum
// ---------------------------------------------------------------------------

// Tab represents which tab is currently active in the detail panel.
// This is a custom type based on int — Go's idiomatic approach to enums.
//
// Using a named type (Tab) instead of a bare int gives us type safety:
// you can't accidentally assign an arbitrary int to a Tab field without
// an explicit conversion, which helps catch bugs at compile time.
//
// See: https://go.dev/ref/spec#Type_definitions
// See: https://go.dev/doc/effective_go#constants
type Tab int

// These constants define the three tabs using iota.
// iota resets to 0 at the start of each const block and increments by 1.
//
//   - TabOverview = 0
//   - TabOutput   = 1
//   - TabConfig   = 2
//
// See: https://go.dev/ref/spec#Iota
const (
	TabOverview Tab = iota // Module info, deps, and status (default tab).
	TabOutput              // Scrollable log of install output.
	TabConfig              // Module-specific configuration options.
)

// tabCount is the total number of tabs. We use this for modular arithmetic
// when cycling through tabs. Keeping it as a constant avoids magic numbers.
const tabCount = 3

// tabNames maps each Tab value to its display label.
// In Go, arrays have a fixed length set at compile time, while slices are
// dynamic. Here we use an array because the size is known and constant.
// See: https://go.dev/tour/moretypes/6
var tabNames = [tabCount]string{"Overview", "Output", "Configuration"}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// TabChangedMsg is sent to the parent model whenever the user switches tabs.
// In Bubble Tea, you communicate between components by defining message types
// (any Go type) and returning them as tea.Cmd functions from Update().
//
// The parent can listen for this message in its own Update() and react
// accordingly (e.g., load data for the newly active tab).
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg
type TabChangedMsg struct {
	Tab Tab // The newly active tab.
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// Model is the tab container for the detail panel. It owns three sub-models
// (one per tab) and delegates messages and rendering to the active one.
//
// In Go, struct fields that start with a lowercase letter are unexported
// (private to the package). This encapsulates internal state so other
// packages can only interact with the model through its public methods.
//
// See: https://go.dev/doc/effective_go#names
// See: https://go.dev/ref/spec#Exported_identifiers
type Model struct {
	// activeTab tracks which tab is currently displayed.
	activeTab Tab

	// width and height define the available space for the detail panel.
	// These are set at construction and updated on terminal resize.
	width  int
	height int

	// Sub-models — one per tab. Each is a value type (not a pointer),
	// so updating them in Update() requires reassigning them back to
	// the struct field. This is the standard Bubble Tea pattern.
	overviewModel OverviewModel
	outputModel   OutputModel
	configModel   ConfigModel

	// moduleName stores the currently displayed module's name.
	// This is used to show which module the detail panel is about.
	moduleName string
}

// NewModel creates and returns an initialised detail panel Model.
//
// In Go, there are no constructors. By convention, "New" functions serve the
// same purpose — they create and return an initialised value. We pass width
// and height so child models know how much space they have to render in.
//
// The subtracted height (3) accounts for the tab bar (1 line of text + 1 line
// of border decoration) and the separator line below it.
//
// See: https://go.dev/doc/effective_go#composite_literals
func NewModel(width, height int) Model {
	// contentHeight is the space available below the tab bar and separator.
	contentHeight := height - 3

	return Model{
		activeTab:     TabOverview,
		width:         width,
		height:        height,
		overviewModel: NewOverviewModel(width, contentHeight),
		outputModel:   NewOutputModel(width, contentHeight),
		configModel:   NewConfigModel(width, contentHeight),
	}
}

// Init is called once when the program starts. For this sub-model, there are
// no initial commands to run, so we return nil.
//
// In Bubble Tea, returning nil from Init() means "do nothing" — the runtime
// will simply wait for the first user event.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is called every time a message (event) arrives. It routes messages
// to the correct sub-model based on which tab is active.
//
// # The Type Switch
//
// Go uses type switches to inspect the concrete type inside an interface.
// msg.(type) checks what kind of message we received. This is the primary
// way Bubble Tea apps handle different event types.
//
// # Why Return (Model, tea.Cmd) Instead of (tea.Model, tea.Cmd)?
//
// Sub-models return their concrete type (Model) rather than the tea.Model
// interface. This avoids the need for type assertions in the parent model
// and is the idiomatic Bubble Tea pattern for nested models.
//
// See: https://go.dev/doc/effective_go#type_switch
// See: https://go.dev/tour/methods/16
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// First, check if the message is a key press for tab switching.
	// We handle this at the container level because tab navigation is the
	// container's responsibility, not the individual tab's.
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// msg.String() returns a human-readable key name (e.g., "tab", "a", "ctrl+c").
		// See: https://pkg.go.dev/charm.land/bubbletea/v2#KeyPressMsg
		if msg.String() == "tab" {
			// Cycle to the next tab using modular arithmetic.
			// The % operator gives the remainder of division — so when we
			// go past the last tab, we wrap back to 0 (TabOverview).
			// See: https://go.dev/ref/spec#Arithmetic_operators
			m.activeTab = (m.activeTab + 1) % Tab(tabCount)

			// Return a command that produces a TabChangedMsg.
			// In Bubble Tea, a tea.Cmd is a function that returns a tea.Msg.
			// The runtime calls it asynchronously and feeds the result back
			// into Update(). We use this to notify the parent model.
			//
			// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
			cmd := func() tea.Msg {
				return TabChangedMsg{Tab: m.activeTab}
			}
			return m, cmd
		}
	}

	// Delegate all other messages to the active tab's sub-model.
	// We capture the returned command so the runtime can execute it.
	var cmd tea.Cmd

	// A switch on a non-interface value doesn't need a type assertion.
	// See: https://go.dev/tour/flowcontrol/9
	switch m.activeTab {
	case TabOverview:
		m.overviewModel, cmd = m.overviewModel.Update(msg)
	case TabOutput:
		m.outputModel, cmd = m.outputModel.Update(msg)
	case TabConfig:
		m.configModel, cmd = m.configModel.Update(msg)
	}

	return m, cmd
}

// View renders the complete detail panel: tab bar + separator + active tab content.
//
// This is a pure function — it reads from the model but never modifies it.
// All state changes happen in Update(). This separation is a core TEA principle.
//
// # String Rendering with Lip Gloss
//
// We build the UI by rendering styled strings and joining them together.
// lipgloss.JoinHorizontal puts strings side by side; lipgloss.JoinVertical
// stacks them top to bottom.
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinHorizontal
// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinVertical
func (m Model) View() string {
	// --- Tab Bar ---
	// Build each tab label, highlighting the active one.
	// We use a slice literal and a for-range loop to iterate over tab names.
	// See: https://go.dev/tour/moretypes/16
	renderedTabs := make([]string, 0, tabCount)
	for i, name := range tabNames {
		// Tab(i) converts int to our Tab type. This is an explicit type
		// conversion, required because Go treats Tab and int as different types.
		// See: https://go.dev/ref/spec#Conversions
		if Tab(i) == m.activeTab {
			// theme.ActiveTab renders with bold pink text and a pink underline.
			renderedTabs = append(renderedTabs, theme.ActiveTab.Render(name))
		} else {
			// theme.InactiveTab renders with dim text and a subtle underline.
			renderedTabs = append(renderedTabs, theme.InactiveTab.Render(name))
		}
	}

	// Join tab labels horizontally, aligned at their bottom edges so the
	// underline decorations line up neatly.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinHorizontal
	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)

	// --- Separator Line ---
	// strings.Repeat creates a thin horizontal rule using a Unicode box-drawing char.
	// We style it with the surface colour so it's visible but not distracting.
	// See: https://pkg.go.dev/strings#Repeat
	separator := lipgloss.NewStyle().
		Foreground(theme.ColorSurface).
		Render(strings.Repeat("─", m.width))

	// --- Tab Content ---
	// Render the active tab's content. Each sub-model's View() returns a string.
	var content string
	switch m.activeTab {
	case TabOverview:
		content = m.overviewModel.View()
	case TabOutput:
		content = m.outputModel.View()
	case TabConfig:
		content = m.configModel.View()
	}

	// Stack the tab bar, separator, and content vertically.
	// lipgloss.Left aligns everything to the left edge.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinVertical
	return lipgloss.JoinVertical(lipgloss.Left, tabBar, separator, content)
}

// ---------------------------------------------------------------------------
// Public Accessors
// ---------------------------------------------------------------------------
// These methods let the parent model interact with the detail panel's state.
// In Go, getter methods are conventionally named without the "Get" prefix
// (e.g., ActiveTab() instead of GetActiveTab()).
// See: https://go.dev/doc/effective_go#Getters

// ActiveTab returns the currently active tab.
func (m Model) ActiveTab() Tab {
	return m.activeTab
}

// SetActiveTab switches to the given tab.
func (m *Model) SetActiveTab(tab Tab) {
	m.activeTab = tab
}

// OverviewModel returns a pointer to the overview sub-model so the parent
// can call methods on it (e.g., SetModule). Returning a pointer allows
// modifications to affect the original, not a copy.
// See: https://go.dev/tour/moretypes/1
func (m *Model) OverviewModel() *OverviewModel {
	return &m.overviewModel
}

// OutputModel returns a pointer to the output sub-model.
func (m *Model) OutputModel() *OutputModel {
	return &m.outputModel
}

// ConfigModel returns a pointer to the config sub-model.
func (m *Model) ConfigModel() *ConfigModel {
	return &m.configModel
}

// SetSize updates the panel dimensions (e.g., on terminal resize).
// It propagates the new size to all child models.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	contentHeight := height - 3
	m.overviewModel.SetSize(width, contentHeight)
	m.outputModel.SetSize(width, contentHeight)
	m.configModel.SetSize(width, contentHeight)
}
