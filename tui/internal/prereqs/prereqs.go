// Package prereqs — prereqs.go is the Bubble Tea sub-model for the prerequisites screen.
//
// This model implements the prerequisites checking UI: it runs checks for
// each required tool asynchronously, displays results with status indicators,
// and allows the user to install missing prerequisites.
//
// # Bubble Tea Sub-Models
//
// In Bubble Tea, complex UIs are built by composing smaller models ("sub-models").
// Each sub-model has its own Init/Update/View cycle, and the parent model
// delegates messages to it. This keeps each model focused and testable.
//
// A sub-model doesn't need to implement tea.Model exactly — it just needs
// methods the parent can call. The parent's Update/View calls the sub-model's
// Update/View and wires everything together.
//
// # Custom Messages (tea.Msg)
//
// In Bubble Tea, all communication happens through messages (tea.Msg).
// You define custom message types as plain Go structs, then handle them in
// Update() with a type switch. This is how async operations report results:
//
//  1. Init() or Update() returns a tea.Cmd (a function that does I/O)
//  2. The Cmd runs in a goroutine and returns a tea.Msg
//  3. Bubble Tea delivers that Msg to Update()
//  4. Update() pattern-matches on the Msg type and updates state
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg
//
// For The Elm Architecture overview: https://guide.elm-lang.org/architecture/
// For Go interfaces: https://go.dev/doc/effective_go#interfaces
package prereqs

import (
	"fmt"
	"os/exec"
	"strings"

	// The Bubbles library provides pre-built UI components for Bubble Tea.
	// We use the help component for displaying keybindings and the key
	// package for defining and matching keyboard shortcuts.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/help
	// See: https://pkg.go.dev/charm.land/bubbles/v2/key
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"

	// The spinner component provides animated loading indicators.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner
	"charm.land/bubbles/v2/spinner"

	// Bubble Tea (aliased as "tea") is the TUI framework.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// Lip Gloss provides CSS-like styling for terminal output.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Custom Types: PrereqStatus Enum
// ---------------------------------------------------------------------------

// PrereqStatus represents the current check status of a single prerequisite.
//
// Go doesn't have a built-in enum keyword. Instead, we define a named type
// based on int and use the iota keyword inside a const block to auto-generate
// sequential values. This is Go's idiomatic way to create enumerations.
//
// See: https://go.dev/ref/spec#Iota
// See: https://go.dev/wiki/Iota
type PrereqStatus int

const (
	// StatusChecking means the prerequisite is still being checked.
	// iota starts at 0 and increments for each constant in the block.
	StatusChecking PrereqStatus = iota

	// StatusOK means the prerequisite is installed (iota = 1).
	StatusOK

	// StatusMissing means the prerequisite is not installed (iota = 2).
	StatusMissing
)

// ---------------------------------------------------------------------------
// Custom Messages (tea.Msg)
// ---------------------------------------------------------------------------
// In Bubble Tea, all events are represented as messages. You define custom
// message types as simple structs and return them from tea.Cmd functions.
// The Bubble Tea runtime delivers these messages to your Update() function.
//
// tea.Msg is defined as an empty interface (interface{}), so ANY Go type
// can be used as a message. We use structs for clarity and to carry data.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Msg
// See: https://go.dev/tour/methods/14 (empty interface)

// PrereqCheckMsg is sent when a single prerequisite check completes.
// Each running check command sends one of these back to Update().
type PrereqCheckMsg struct {
	// Name identifies which prerequisite was checked.
	Name string

	// Installed is true if the prerequisite was found on the system.
	Installed bool
}

// AllChecksCompleteMsg is sent when every prerequisite has been checked.
// This is an internal signal to transition from the "checking" state
// to the "results" state. An empty struct (struct{}) carries no data —
// its mere presence IS the message.
//
// See: https://go.dev/ref/spec#Size_and_alignment_guarantees (empty struct = zero bytes)
type AllChecksCompleteMsg struct{}

// PrereqsPassedMsg is the exported result message sent to the parent model.
// When all prerequisites are satisfied and the user presses Enter, this
// message bubbles up to the root app model to advance to the next screen.
//
// Exported types (capitalised names) are visible outside the package.
// The parent app model type-switches on this to know prereqs are done.
// See: https://go.dev/doc/effective_go#names
type PrereqsPassedMsg struct{}

// InstallCompleteMsg is sent when the installation command finishes.
// It carries an error field: nil on success, non-nil on failure.
type InstallCompleteMsg struct {
	// Err is nil if installation succeeded, or an error describing what went wrong.
	// In Go, the built-in error interface has a single method: Error() string.
	// See: https://go.dev/doc/effective_go#errors
	Err error
}

// ---------------------------------------------------------------------------
// Key Bindings
// ---------------------------------------------------------------------------
// We define a local keyMap type for this screen's keyboard shortcuts.
// This implements the help.KeyMap interface so the Bubbles help component
// can automatically render a help bar from our bindings.
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/help#KeyMap
// See: https://pkg.go.dev/charm.land/bubbles/v2/key#Binding

// prereqKeyMap defines the keyboard shortcuts available on the prerequisites screen.
type prereqKeyMap struct {
	Install key.Binding
	Proceed key.Binding
	Quit    key.Binding
}

// defaultPrereqKeyMap returns the default key bindings for this screen.
//
// key.NewBinding creates a binding with:
//   - key.WithKeys(): The actual key(s) that trigger it (matched against tea.KeyPressMsg)
//   - key.WithHelp(): A short key label + description for the help bar
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/key#NewBinding
var defaultPrereqKeyMap = prereqKeyMap{
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install missing"),
	),
	Proceed: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp returns bindings for the compact (one-line) help view.
// This satisfies the help.KeyMap interface — Go interfaces are satisfied
// implicitly (no "implements" keyword needed).
// See: https://go.dev/doc/effective_go#interfaces
func (k prereqKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Proceed, k.Quit}
}

// FullHelp returns bindings grouped by category for the expanded help view.
// The outer slice represents columns; inner slices are bindings in each column.
func (k prereqKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Install, k.Proceed, k.Quit},
	}
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// Model is the Bubble Tea sub-model for the prerequisites checking screen.
//
// It holds all the state needed to:
//  1. Track which prereqs are being checked / installed / missing
//  2. Animate a spinner while checks run
//  3. Display results and accept user input
//
// In Go, struct fields that start with a lowercase letter are unexported
// (private to this package). This encapsulation prevents other packages
// from directly manipulating internal state — they must use exported methods.
//
// See: https://go.dev/doc/effective_go#names
// See: https://go.dev/ref/spec#Struct_types
type Model struct {
	// prereqs is the full list of prerequisites to check.
	prereqs []Prereq

	// statuses tracks the check result for each prerequisite by name.
	// Maps in Go are reference types — they're like hash tables / dictionaries.
	// See: https://go.dev/doc/effective_go#maps
	statuses map[string]PrereqStatus

	// checking is true while prerequisite checks are still running.
	checking bool

	// installing is true while apt install is running.
	installing bool

	// allOK is true when every prerequisite is installed.
	allOK bool

	// installErr holds the error from a failed installation attempt.
	// It's nil when no error has occurred.
	installErr error

	// checkedCount tracks how many individual checks have reported back.
	// When this equals len(prereqs), all checks are complete.
	checkedCount int

	// spinner is a Bubbles spinner sub-model for the loading animation.
	// Spinner.View() returns the current animation frame as a string.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner
	spinner spinner.Model

	// keys holds the keyboard shortcut definitions for this screen.
	keys prereqKeyMap

	// help is a Bubbles help sub-model that renders the keybinding help bar.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/help
	help help.Model
}

// New creates and returns an initialised prerequisites model.
//
// In Go, constructors are just regular functions — conventionally named
// New or New<Type>. They return an initialised value (not a pointer, unless
// the struct is large or needs to be shared). Go has no "new" keyword for
// custom types; instead you use struct literals: Type{field: value}.
//
// See: https://go.dev/doc/effective_go#composite_literals
// See: https://go.dev/doc/effective_go#allocation_new
func New() Model {
	prereqs := GetRequiredPrereqs()

	// make(map[K]V) creates an empty map. Maps must be initialised with make()
	// before use — a nil map will panic on write (but reads return zero values).
	// See: https://go.dev/doc/effective_go#maps
	statuses := make(map[string]PrereqStatus, len(prereqs))

	// Initialise every prereq's status to "checking" (the default state).
	for _, p := range prereqs {
		statuses[p.Name] = StatusChecking
	}

	// Create the spinner using the functional options pattern.
	// spinner.New() accepts Option functions that configure the model.
	// spinner.WithSpinner sets the animation frames (Dot = braille dots).
	// spinner.WithStyle sets the Lip Gloss style (colour, bold, etc.).
	//
	// The functional options pattern is a common Go idiom for configurable
	// constructors. Each option is a function that modifies the struct.
	// See: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	// See: https://pkg.go.dev/charm.land/bubbles/v2/spinner#New
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(theme.ColorCyan)),
	)

	return Model{
		prereqs:  prereqs,
		statuses: statuses,
		checking: true,
		spinner:  s,
		keys:     defaultPrereqKeyMap,
		help:     help.New(),
	}
}

// ---------------------------------------------------------------------------
// Init — The Entry Point
// ---------------------------------------------------------------------------

// Init returns the initial command(s) to run when this model starts.
//
// # How tea.Cmd Works
//
// A tea.Cmd is a function with the signature: func() tea.Msg
// Bubble Tea runs each Cmd in a separate goroutine. When the Cmd returns
// a tea.Msg, that message is delivered to Update(). This is how you perform
// async I/O (network calls, shell commands, timers) without blocking the UI.
//
// tea.Batch() combines multiple Cmds into one. All batched commands run
// concurrently — their result messages arrive in whatever order they finish.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Batch
func (m Model) Init() tea.Cmd {
	// Build a slice of commands: one spinner tick + one check per prereq.
	//
	// []tea.Cmd{...} is a slice literal. We start with the spinner's Tick
	// command to begin the loading animation.
	//
	// m.spinner.Tick is a method value — in Go, you can reference a method
	// without calling it to get a function value. Since Tick() has the
	// signature func() tea.Msg, it matches the tea.Cmd type exactly.
	// See: https://go.dev/ref/spec#Method_values
	cmds := []tea.Cmd{m.spinner.Tick}

	// Create one Cmd for each prerequisite check. Each runs concurrently
	// and sends back a PrereqCheckMsg when done.
	for _, p := range m.prereqs {
		// IMPORTANT: capture the loop variable in a local copy.
		// Go's range loop reuses the same variable address on each iteration,
		// so if we closed over 'p' directly, all goroutines would see the
		// LAST value of p. Creating a local copy fixes this.
		// See: https://go.dev/doc/faq#closures_and_goroutines
		//
		// Note: As of Go 1.22+, the loop variable is per-iteration, so this
		// is technically no longer needed. We keep it for clarity and for
		// compatibility with older Go versions.
		p := p
		cmds = append(cmds, checkPrereqCmd(p))
	}

	// tea.Batch combines all commands into one. Bubble Tea runs them
	// concurrently and delivers each result message to Update().
	return tea.Batch(cmds...)
}

// checkPrereqCmd creates a tea.Cmd that checks a single prerequisite.
//
// This is a "command factory" — a function that returns a tea.Cmd. The
// returned Cmd is a closure that captures the Prereq value and runs
// CheckPrereq() when executed by the Bubble Tea runtime.
//
// # Closures in Go
//
// A closure is a function that references variables from its enclosing scope.
// Here, the anonymous function func() tea.Msg "closes over" the p variable,
// meaning it retains access to p even after checkPrereqCmd returns.
//
// See: https://go.dev/tour/moretypes/25 (function closures)
// See: https://go.dev/doc/effective_go#functions
func checkPrereqCmd(p Prereq) tea.Cmd {
	return func() tea.Msg {
		installed := CheckPrereq(p)
		return PrereqCheckMsg{
			Name:      p.Name,
			Installed: installed,
		}
	}
}

// installCmd creates a tea.Cmd that runs the installation commands.
//
// The command runs "sudo apt update && sudo apt install -y <packages>"
// in a shell and returns an InstallCompleteMsg with the result.
//
// See: https://pkg.go.dev/os/exec#Command
func installCmd(missing []Prereq) tea.Cmd {
	return func() tea.Msg {
		commands := InstallPrereqs(missing)
		// strings.Join concatenates the commands with " && " to run them
		// sequentially in a single shell invocation.
		// See: https://pkg.go.dev/strings#Join
		combined := strings.Join(commands, " && ")

		// #nosec G204 — command is built from controlled internal data
		cmd := exec.Command("sh", "-c", combined)
		err := cmd.Run()

		return InstallCompleteMsg{Err: err}
	}
}

// ---------------------------------------------------------------------------
// Update — Message Handling
// ---------------------------------------------------------------------------

// Update processes incoming messages and returns the updated model + any new commands.
//
// This is where all state changes happen. Bubble Tea calls Update() every time
// a message arrives (key press, timer tick, async operation result, etc.).
//
// # Type Switch Pattern
//
// Go's type switch checks the concrete type of an interface value.
// Since tea.Msg is an empty interface, we use msg.(type) to determine
// which specific message type we received and handle it accordingly.
//
// See: https://go.dev/tour/methods/16 (type assertions)
// See: https://go.dev/doc/effective_go#type_switch
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// -----------------------------------------------------------------------
	// PrereqCheckMsg — A single prerequisite check has completed
	// -----------------------------------------------------------------------
	case PrereqCheckMsg:
		// Update the status map for this specific prerequisite.
		if msg.Installed {
			m.statuses[msg.Name] = StatusOK
		} else {
			m.statuses[msg.Name] = StatusMissing
		}

		// Increment the counter of completed checks.
		m.checkedCount++

		// When all checks have reported back, determine the overall result.
		if m.checkedCount >= len(m.prereqs) {
			m.checking = false

			// Assume all OK, then check for any missing prerequisites.
			// This is a common Go pattern: optimistic default, then invalidate.
			m.allOK = true
			for _, s := range m.statuses {
				if s == StatusMissing {
					m.allOK = false
					// break exits the innermost for loop early — no need to
					// keep checking once we know something is missing.
					// See: https://go.dev/ref/spec#Break_statements
					break
				}
			}

			// Enable/disable keybindings based on the result.
			// key.SetEnabled toggles whether a binding responds to key.Matches.
			// See: https://pkg.go.dev/charm.land/bubbles/v2/key#SetEnabled
			m.keys.Install.SetEnabled(!m.allOK)
			m.keys.Proceed.SetEnabled(m.allOK)
		}

		return m, nil

	// -----------------------------------------------------------------------
	// InstallCompleteMsg — The installation command has finished
	// -----------------------------------------------------------------------
	case InstallCompleteMsg:
		m.installing = false

		if msg.Err != nil {
			// Store the error so View() can display it.
			m.installErr = msg.Err
			return m, nil
		}

		// Installation succeeded — recheck all prerequisites.
		// Reset state and re-run Init() to start fresh checks.
		m.checking = true
		m.checkedCount = 0
		m.allOK = false
		m.installErr = nil
		for _, p := range m.prereqs {
			m.statuses[p.Name] = StatusChecking
		}

		return m, m.Init()

	// -----------------------------------------------------------------------
	// Key Presses — User interaction
	// -----------------------------------------------------------------------
	// tea.KeyPressMsg is sent every time the user presses a key.
	// In Bubble Tea v2, this replaces the v1 tea.KeyMsg type.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#KeyPressMsg
	case tea.KeyPressMsg:
		// key.Matches checks if the pressed key matches one of our bindings.
		// It's a generic function: key.Matches[Key fmt.Stringer](k, bindings...).
		// tea.KeyPressMsg implements fmt.Stringer, so it works directly.
		// See: https://pkg.go.dev/charm.land/bubbles/v2/key#Matches
		switch {
		case key.Matches(msg, m.keys.Quit):
			// tea.Quit is a special Cmd that tells the Bubble Tea runtime
			// to shut down the program gracefully.
			// See: https://pkg.go.dev/charm.land/bubbletea/v2#Quit
			return m, tea.Quit

		case key.Matches(msg, m.keys.Install):
			// Only allow installing when checks are done and something is missing.
			if !m.checking && !m.allOK && !m.installing {
				m.installing = true
				m.installErr = nil
				return m, installCmd(m.getMissing())
			}

		case key.Matches(msg, m.keys.Proceed):
			// Only allow proceeding when all prerequisites are installed.
			if !m.checking && m.allOK {
				// Return a Cmd that sends PrereqsPassedMsg to the parent.
				// This anonymous function is a tea.Cmd (func() tea.Msg).
				// The parent model will receive PrereqsPassedMsg in its Update().
				return m, func() tea.Msg { return PrereqsPassedMsg{} }
			}
		}

	// -----------------------------------------------------------------------
	// Spinner tick messages — Animation updates
	// -----------------------------------------------------------------------
	// The spinner.TickMsg is handled by the spinner sub-model.
	// We delegate to spinner.Update() which advances the animation frame
	// and returns the next tick Cmd to continue the animation loop.
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// If no handler matched, return the model unchanged with no command.
	return m, nil
}

// getMissing returns only the prerequisites that are not installed.
//
// This is a method with a value receiver (m Model), meaning it receives a
// copy of the Model. Value receivers are used when the method doesn't need
// to modify the original struct.
//
// See: https://go.dev/tour/methods/4 (pointer vs value receivers)
func (m Model) getMissing() []Prereq {
	// var declares a variable. A nil slice is valid and can be appended to.
	// See: https://go.dev/ref/spec#Variable_declarations
	var missing []Prereq

	for _, p := range m.prereqs {
		if m.statuses[p.Name] == StatusMissing {
			// append() is a built-in that adds elements to a slice.
			// It may allocate a new underlying array if capacity is exceeded.
			// See: https://pkg.go.dev/builtin#append
			missing = append(missing, p)
		}
	}

	return missing
}

// ---------------------------------------------------------------------------
// View — Rendering
// ---------------------------------------------------------------------------

// View renders the prerequisites screen as a styled string.
//
// IMPORTANT: View() must be a pure function. It reads from the model but
// never modifies it or performs I/O. All side effects happen in Update()
// through Cmds. This is a core principle of The Elm Architecture.
//
// We use strings.Builder for efficient string concatenation. Unlike using
// the + operator repeatedly (which creates a new string each time),
// strings.Builder uses an internal buffer and builds the final string once.
//
// See: https://pkg.go.dev/strings#Builder
// See: https://go.dev/doc/effective_go#printing
func (m Model) View() string {
	// strings.Builder is Go's efficient way to build strings incrementally.
	// It implements the io.Writer interface, so it works with fmt.Fprintf too.
	// See: https://pkg.go.dev/strings#Builder
	var b strings.Builder

	// -- Title --
	b.WriteString(theme.Title.Render("🔍 Pre-requisite Checks"))
	b.WriteString("\n\n")

	// -- Column headers --
	// We use fmt.Sprintf to format a fixed-width table row.
	// The %-Ns format verb left-aligns a string in a field of width N.
	// See: https://pkg.go.dev/fmt#hdr-Printing
	header := fmt.Sprintf(
		"  %s  %-20s %-35s %s",
		theme.Subtitle.Render("•"),
		theme.Subtitle.Render("Name"),
		theme.Subtitle.Render("Description"),
		theme.Subtitle.Render("Docs"),
	)
	b.WriteString(header)
	b.WriteString("\n")

	// A separator line for visual clarity.
	b.WriteString(theme.DimText.Render("  "+strings.Repeat("─", 80)))
	b.WriteString("\n")

	// -- Prerequisite rows --
	// Iterate over each prerequisite and render its status row.
	for _, p := range m.prereqs {
		status := m.statuses[p.Name]

		// Determine the status icon based on the check result.
		// The spinner.View() returns the current animation frame when checking.
		var statusIcon string
		switch status {
		case StatusChecking:
			// Show the animated spinner character while still checking.
			statusIcon = m.spinner.View()
		case StatusOK:
			// theme.StatusInstalled has SetString("✓") — Render() outputs "✓"
			// with green foreground colour applied.
			statusIcon = theme.StatusInstalled.Render()
		case StatusMissing:
			// theme.StatusNotInstalled has SetString("✗") — Render() outputs "✗"
			// with red foreground colour applied.
			statusIcon = theme.StatusNotInstalled.Render()
		}

		// Render each column with appropriate styling.
		name := theme.NormalText.Render(p.Name)
		desc := theme.DimText.Render(p.Description)
		url := theme.URLStyle.Render(p.DocsURL)

		// Format the row with fixed-width columns for alignment.
		// Note: ANSI escape codes (from Lip Gloss styling) add invisible
		// characters, so the visual alignment may not be pixel-perfect.
		// For a production app, you'd use lipgloss.Width() to measure
		// visible width and pad accordingly.
		row := fmt.Sprintf("  %s  %-20s %-35s %s", statusIcon, name, desc, url)
		b.WriteString(row)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// -- Status message --
	// Show different messages depending on the current state.
	switch {
	case m.installing:
		// Show a spinner + message while installing.
		b.WriteString(
			m.spinner.View() + " " +
				theme.WarningText.Render("Installing prerequisites... This may take a moment."),
		)
		b.WriteString("\n")

	case m.checking:
		// Show a message while checks are still running.
		b.WriteString(
			m.spinner.View() + " " +
				theme.DimText.Render("Checking prerequisites..."),
		)
		b.WriteString("\n")

	case m.installErr != nil:
		// Show the installation error.
		// fmt.Sprintf formats a string using Printf-style verbs.
		// %v prints the default representation of any value.
		// See: https://pkg.go.dev/fmt#hdr-Printing
		b.WriteString(theme.ErrorText.Render(
			fmt.Sprintf("✗ Installation failed: %v", m.installErr),
		))
		b.WriteString("\n")
		b.WriteString(theme.DimText.Render("  You may need to run the install command manually with sudo."))
		b.WriteString("\n")

	case m.allOK:
		// All prerequisites are installed — show success message.
		b.WriteString(theme.SuccessText.Render("✓ All prerequisites met! Press Enter to continue."))
		b.WriteString("\n")

	default:
		// Some prerequisites are missing — list them and offer to install.
		missing := m.getMissing()

		b.WriteString(theme.ErrorText.Render("✗ Missing prerequisites:"))
		b.WriteString("\n")

		for _, p := range missing {
			b.WriteString(
				"    " +
					theme.ErrorText.Render("• ") +
					theme.NormalText.Render(p.Name) +
					theme.DimText.Render(" — "+p.Description),
			)
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(theme.WarningText.Render("Press 'i' to install missing prerequisites"))
		b.WriteString("\n")
	}

	// -- Help bar --
	// The Bubbles help component automatically renders our key bindings.
	// help.View() accepts any type that satisfies the help.KeyMap interface
	// (i.e., has ShortHelp() and FullHelp() methods).
	// See: https://pkg.go.dev/charm.land/bubbles/v2/help#Model.View
	b.WriteString("\n")
	b.WriteString(theme.HelpStyle.Render(m.help.View(m.keys)))

	// b.String() returns the accumulated string from all WriteString calls.
	return b.String()
}
