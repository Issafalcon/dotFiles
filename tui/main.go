// Package main is the entry point for the DotFiles TUI application.
//
// This application provides an interactive terminal user interface for managing
// dotfile configurations using GNU Stow. It's built using the Charm ecosystem:
//
//   - Bubble Tea: The TUI framework based on The Elm Architecture
//     https://pkg.go.dev/charm.land/bubbletea/v2
//   - Lip Gloss: CSS-like styling for terminal output
//     https://pkg.go.dev/charm.land/lipgloss/v2
//   - Bubbles: Pre-built UI components (lists, spinners, text inputs, etc.)
//     https://pkg.go.dev/charm.land/bubbles/v2
//
// # The Elm Architecture (TEA)
//
// Bubble Tea apps follow The Elm Architecture pattern:
//
//  1. Model: A struct holding all application state
//  2. Init(): Returns an initial command to run (e.g., fetch data)
//  3. Update(msg): Receives messages (events) and returns updated model + commands
//  4. View(): Renders the current state as a string for display
//
// Messages flow in one direction: Event → Update → Model → View
// This makes the app predictable and easy to reason about.
//
// For more on The Elm Architecture, see: https://guide.elm-lang.org/architecture/
// For Go basics, see: https://go.dev/doc/
package main

import (
	"fmt"
	"os"

	// tea is the conventional alias for the Bubble Tea framework.
	// In Go, you can alias imports to shorter names for convenience.
	// See: https://go.dev/ref/spec#Import_declarations
	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/dotfiles-tui/internal/app"
)

// version is set at build time via -ldflags "-X main.version=v1.2.3".
// If not set (e.g., during `go run .`), it defaults to "dev".
// See: https://pkg.go.dev/cmd/link
var version = "dev"

func main() {
	// Handle --version flag for quick version checks.
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("dotfiles-tui %s\n", version)
		os.Exit(0)
	}

	// Create the root application model.
	// In Go, short variable declaration (:=) infers the type automatically.
	// See: https://go.dev/tour/basics/10
	initialModel := app.NewModel()

	// tea.NewProgram creates a new Bubble Tea program.
	// In Bubble Tea v2, features like alternate screen and mouse mode are set
	// declaratively in the View() method rather than as program options.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#NewProgram
	p := tea.NewProgram(initialModel)

	// Send the program reference to the model so it can use p.Send()
	// from background goroutines (e.g., streaming install output).
	// This runs in a goroutine because p.Send() blocks until the program
	// is ready to receive messages (which happens after p.Run() starts).
	go func() {
		p.Send(app.ProgramReadyMsg{Program: p})
	}()

	// p.Run() starts the event loop. It blocks until the program exits.
	// The underscore (_) discards the final model — we don't need it after exit.
	// In Go, you must explicitly handle or discard return values.
	// See: https://go.dev/doc/effective_go#blank
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running dotfiles TUI: %v\n", err)
		os.Exit(1)
	}
}
