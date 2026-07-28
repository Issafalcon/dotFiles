# DotFiles TUI — Developer Guide

## Overview

The DotFiles TUI is an interactive terminal application for installing, stowing, and uninstalling modules under `modules/`. Install/uninstall logic lives in each module's shell scripts; Go holds metadata only. It's built using the [Charm](https://charm.sh/) ecosystem.

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
/
├── go.mod                      # Go module (repo root)
├── modules/<name>/             # Per-module: module.yaml + install.sh + configs
├── internal/                   # TUI core packages
│   ├── app/
│   ├── config/                 # ~/.config/dotfiles-tui + modules path resolution
│   ├── docs/                   # Embedded ADDING_MODULES / TROUBLESHOOTING
│   ├── module/                 # Registry + YAML loader (runtime discovery)
│   ├── sidebar/                # List + search + category filter
│   └── ...
└── tui/
    ├── main.go
    ├── Makefile                # build / run / appimage (no modules in AppImage)
    ├── appimage/               # AppRun + desktop/icon
    └── docs/
```

Modules are discovered at runtime from `module.yaml` under the configured modules
directory (env `DOTFILES_MODULES_DIR`, config, or auto-detected repo `modules/`).
The AppImage ships the TUI + docs only — zero module packages.

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
