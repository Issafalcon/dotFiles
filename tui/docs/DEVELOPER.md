# DotFiles TUI — Developer Guide

## Overview

The DotFiles TUI is an interactive terminal application that replaces the shell-based `bootstrap.sh` workflow with a full-featured module management interface. It's built in Go using the [Charm](https://charm.sh/) ecosystem.

## Prerequisites

- **Go 1.25+** — [Install Go](https://go.dev/dl/)
- **GNU Stow** — `sudo apt install stow` (used for symlink management)
- **A Nerd Font** — [Nerd Fonts](https://www.nerdfonts.com/) for icons to render properly

## Quick Start

```bash
# Clone the dotfiles repo (if you haven't already)
git clone https://github.com/issafalcon/dotfiles.git
cd dotfiles/tui

# Run the app directly
make run

# Or build a binary
make build
./build/dotfiles-tui

# Install system-wide
make install
dotfiles-tui
```

## Project Structure

```
tui/
├── main.go                     # Entry point
├── go.mod / go.sum             # Go module dependencies
├── Makefile                    # Build targets
├── internal/                   # Private packages (Go convention)
│   ├── app/                    # Root application model
│   │   ├── app.go              # Main model (Init/Update/View)
│   │   └── keys.go             # Keyboard shortcuts
│   ├── theme/                  # Visual styling
│   │   └── theme.go            # Colours, styles, icons
│   ├── prereqs/                # Prerequisites checking
│   │   ├── checker.go          # Detection logic
│   │   └── prereqs.go          # TUI screen
│   ├── sidebar/                # Left panel — module list
│   │   ├── sidebar.go          # List component
│   │   └── filter.go           # Fuzzy search
│   ├── detail/                 # Right panel — module details
│   │   ├── detail.go           # Tab container
│   │   ├── overview.go         # Overview tab
│   │   ├── output.go           # Install output tab
│   │   └── config.go           # Configuration tab
│   ├── installer/              # Installation engine
│   │   ├── installer.go        # Parallel orchestrator
│   │   ├── runner.go           # Command execution
│   │   └── verify.go           # Post-install checks
│   ├── module/                 # Module definitions
│   │   ├── module.go           # Module struct
│   │   ├── registry.go         # Module registry
│   │   └── modules/            # Individual module defs
│   │       ├── template.go     # New module template
│   │       ├── editors.go      # nvim, vimspector
│   │       ├── shell.go        # zsh, powershell
│   │       ├── languages.go    # go, rust, python, etc.
│   │       └── ...
│   ├── popup/                  # Modal dialogs
│   │   ├── popup.go            # Generic overlay
│   │   ├── confirm.go          # Install confirmation
│   │   ├── input.go            # User input
│   │   └── help.go             # Help overlay
│   └── utils/                  # Shared utilities
│       ├── exec.go             # Shell command runner
│       ├── stow.go             # GNU Stow wrapper
│       ├── browser.go          # URL opener
│       └── detect.go           # Software detection
└── docs/
    ├── DEVELOPER.md            # This file
    └── ADDING_MODULES.md       # Module creation guide
```

## Architecture: The Elm Architecture (TEA)

This app follows [The Elm Architecture](https://guide.elm-lang.org/architecture/), implemented by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

### The Pattern

Every component in the app follows the same three-function pattern:

1. **`Init() tea.Cmd`** — Called once on startup. Returns initial commands (I/O operations).
2. **`Update(msg tea.Msg) (tea.Model, tea.Cmd)`** — Called on every event. Processes the event and returns updated state + optional new commands.
3. **`View() string`** — Called after every Update. Returns the UI as a string. Must be a **pure function** — no side effects.

### Message Flow

```
User Input / Timer / I/O Result
        ↓
    tea.Msg (a message)
        ↓
    Update(msg) → new Model + optional Cmd
        ↓
    View() → rendered string
        ↓
    Terminal Output
```

### Commands (tea.Cmd)

A `tea.Cmd` is a function that performs I/O and returns a `tea.Msg`:

```go
// A command that checks if git is installed
func checkGit() tea.Msg {
    _, err := exec.LookPath("git")
    return PrereqCheckMsg{Name: "git", Installed: err == nil}
}
```

Commands are the **only** way to perform side effects. The Update function returns them, and Bubble Tea runs them asynchronously.

## Libraries Used

| Library | Import Path | Purpose |
|---------|-------------|---------|
| [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) | `charm.land/bubbletea/v2` | TUI framework |
| [Lip Gloss v2](https://github.com/charmbracelet/lipgloss) | `charm.land/lipgloss/v2` | Terminal styling |
| [Bubbles v2](https://github.com/charmbracelet/bubbles) | `charm.land/bubbles/v2` | UI components |
| [Huh v2](https://github.com/charmbracelet/huh) | `charm.land/huh/v2` | Forms & prompts |
| [Glamour](https://github.com/charmbracelet/glamour) | `github.com/charmbracelet/glamour` | Markdown rendering |

## Key Go Concepts Used

### Interfaces (Implicit Satisfaction)

Go interfaces are satisfied implicitly — no `implements` keyword needed:

```go
// Any type with these methods is a tea.Model
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() View
}
```

See: https://go.dev/doc/effective_go#interfaces

### Goroutines & Channels

Used for parallel installations:

```go
go func() {
    result := runCommand(cmd)
    resultChan <- result  // send result to channel
}()
```

See: https://go.dev/tour/concurrency/1

### Type Switches

Used extensively in Update functions:

```go
switch msg := msg.(type) {
case tea.KeyPressMsg:
    // handle key press
case tea.WindowSizeMsg:
    // handle resize
}
```

See: https://go.dev/tour/methods/16

### Struct Embedding (Composition)

Go uses composition instead of inheritance:

```go
type Model struct {
    sidebar  sidebar.Model   // embeds the sidebar sub-model
    detail   detail.Model    // embeds the detail sub-model
}
```

See: https://go.dev/doc/effective_go#embedding

## Testing

```bash
make test
```

## Debugging

Since the TUI controls stdin/stdout, use file-based logging:

```go
import tea "charm.land/bubbletea/v2"

// At program start
f, _ := tea.LogToFile("debug.log", "debug")
defer f.Close()
```

Then in another terminal: `tail -f debug.log`

## Linting

```bash
make lint  # requires golangci-lint
```
