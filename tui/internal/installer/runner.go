// Package installer provides the installation engine for the DotFiles TUI.
//
// It integrates with Bubble Tea's message-passing architecture to run shell
// commands asynchronously while keeping the UI responsive. Commands are executed
// in the background via goroutines, and progress is streamed to the UI via
// Bubble Tea messages sent through the *tea.Program reference.
//
// Key Go concepts used here:
//   - tea.Cmd: A function that returns a tea.Msg (https://pkg.go.dev/charm.land/bubbletea/v2#Cmd)
//   - tea.Msg: An interface{} (any value) that carries information through the Update loop
//   - Closures: Functions that capture variables from their enclosing scope
//   - Goroutines: Lightweight concurrent execution (https://go.dev/doc/effective_go#goroutines)
//   - p.Send(): Injecting messages from outside the Update loop
//
// Architecture:
//
// Install commands run in a background goroutine. Each line of stdout/stderr
// is streamed to the Output pane via p.Send(InstallOutputMsg{...}). When
// all commands finish (or one fails), an InstallCompleteMsg is returned as
// the final tea.Msg from the tea.Cmd. This keeps the TUI fully responsive
// throughout the installation process.
//
// For commands that require sudo, RunSudoAuth() runs "sudo -v" via
// tea.ExecProcess (which briefly suspends the TUI for password entry),
// caching the sudo credential. All subsequent commands then run
// non-interactively.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2
// See: https://github.com/charmbracelet/bubbletea/tree/main/tutorials
package installer

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/dotfiles-tui/internal/utils"
)

// --- Bubble Tea Message Types ---

// InstallStartMsg is sent when a module's installation begins.
type InstallStartMsg struct {
	ModuleName string
}

// InstallOutputMsg carries a single line of output from a running installation.
// The UI displays this in the Output tab's scrollable viewport.
type InstallOutputMsg struct {
	ModuleName string
	Line       string
	IsStderr   bool
}

// InstallCompleteMsg is sent when a module's installation finishes.
type InstallCompleteMsg struct {
	ModuleName string
	Success    bool
	Error      error
}

// UninstallStartMsg is sent when a module's uninstallation begins.
type UninstallStartMsg struct {
	ModuleName string
}

// UninstallCompleteMsg is sent when a module's uninstallation finishes.
// The app uses Success to update the sidebar status and show a result message.
type UninstallCompleteMsg struct {
	ModuleName string
	Success    bool
	Error      error
}

// InstallProgressMsg reports step-by-step progress during multi-command installations.
type InstallProgressMsg struct {
	ModuleName string
	Step       int
	TotalSteps int
}

// SudoAuthCompleteMsg is sent after "sudo -v" finishes (via tea.ExecProcess).
// The Update handler uses this to proceed with the streaming install.
type SudoAuthCompleteMsg struct {
	ModuleName string
	Error      error
}

// --- Sudo helpers ---

// NeedsSudo returns true if any of the given commands contain "sudo".
func NeedsSudo(commands []string) bool {
	for _, cmd := range commands {
		if strings.Contains(cmd, "sudo") {
			return true
		}
	}
	return false
}

// RunSudoAuth runs "sudo -v" via tea.ExecProcess to cache the user's sudo
// credentials. This briefly suspends the TUI so the terminal can display
// the password prompt. After authentication, the TUI resumes and all
// subsequent sudo commands run without prompting (credentials are cached
// for ~15 minutes by default).
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#ExecProcess
func RunSudoAuth(moduleName string) tea.Cmd {
	c := exec.Command("sudo", "-v")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return SudoAuthCompleteMsg{ModuleName: moduleName, Error: err}
	})
}

// --- Streaming install ---

// RunInstallStreaming returns a tea.Cmd that runs all install commands in
// a background goroutine, streaming each line of output to the UI via
// p.Send(). The final message returned by the tea.Cmd is InstallCompleteMsg.
//
// Because p.Send() injects messages from outside the Update loop, the TUI
// remains fully responsive (scrolling, tab switching, etc.) while commands
// run. This replaces the old tea.ExecProcess approach which suspended the
// entire TUI for each command.
//
// Parameters:
//   - p: The running Bubble Tea program, used to send streaming messages.
//   - moduleName: The module being installed.
//   - commands: Shell commands to execute sequentially.
//   - dotfilesDir: The dotfiles repository root (for stow).
//   - stowEnabled: Whether to run stow after all commands succeed.
func RunInstallStreaming(p *tea.Program, moduleName string, commands []string, dotfilesDir string, stowEnabled bool) tea.Cmd {
	return func() tea.Msg {
		p.Send(InstallStartMsg{ModuleName: moduleName})

		for i, cmdStr := range commands {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       fmt.Sprintf("\n▸ Step %d/%d: %s", i+1, len(commands), cmdStr),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			err := utils.RunCommandStreaming(ctx, cmdStr, func(line string, isStderr bool) {
				p.Send(InstallOutputMsg{
					ModuleName: moduleName,
					Line:       line,
					IsStderr:   isStderr,
				})
			})
			cancel()

			if err != nil {
				return InstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("step %d/%d %q: %w", i+1, len(commands), cmdStr, err),
				}
			}
		}

		if stowEnabled {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "\n▸ Creating stow symlinks...",
			})
			if err := utils.Stow(moduleName, dotfilesDir); err != nil {
				return InstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("stow: %w", err),
				}
			}
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "✓ Stow links created",
			})
		}

		if err := utils.SetModuleInstalled(moduleName); err != nil {
			return InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("tracking install: %w", err),
			}
		}

		return InstallCompleteMsg{
			ModuleName: moduleName,
			Success:    true,
		}
	}
}

// RunInstallWithSend executes install commands and sends progress messages via
// a provided send function. This is used by the Orchestrator for real-time
// streaming output to the Bubble Tea UI.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - moduleName: The module being installed.
//   - commands: Shell commands to execute sequentially.
//   - dotfilesDir: The dotfiles repository root.
//   - stowEnabled: Whether to stow after all commands succeed.
//   - send: A function to dispatch tea.Msg values to the Bubble Tea runtime.
func RunInstallWithSend(ctx context.Context, moduleName string, commands []string, dotfilesDir string, stowEnabled bool, send func(tea.Msg)) {
	send(InstallStartMsg{ModuleName: moduleName})

	for i, cmd := range commands {
		send(InstallProgressMsg{
			ModuleName: moduleName,
			Step:       i + 1,
			TotalSteps: len(commands),
		})

		err := utils.RunCommandStreaming(ctx, cmd, func(line string, isStderr bool) {
			send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       line,
				IsStderr:   isStderr,
			})
		})

		if err != nil {
			send(InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("command %d/%d failed: %q: %w", i+1, len(commands), cmd, err),
			})
			return
		}
	}

	if stowEnabled {
		if err := utils.Stow(moduleName, dotfilesDir); err != nil {
			send(InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("stow failed: %w", err),
			})
			return
		}
	}

	if err := utils.SetModuleInstalled(moduleName); err != nil {
		send(InstallCompleteMsg{
			ModuleName: moduleName,
			Success:    false,
			Error:      fmt.Errorf("tracking install: %w", err),
		})
		return
	}

	send(InstallCompleteMsg{
		ModuleName: moduleName,
		Success:    true,
	})
}

// --- Streaming uninstall ---

// RunUninstallStreaming returns a tea.Cmd that runs all uninstall commands in
// a background goroutine, then removes stow symlinks and updates the tracking
// file. Output is streamed line-by-line to the Output pane via p.Send().
//
// The flow mirrors RunInstallStreaming but in reverse:
//  1. Run each UninstallCommand sequentially (e.g., "sudo apt-get remove -y nvim").
//  2. If StowEnabled, run Unstow to remove config symlinks.
//  3. Mark the module as uninstalled in the tracking file.
//
// If any uninstall command fails, the process stops and reports the error,
// but the module is NOT marked as uninstalled (it may be partially removed).
//
// Parameters:
//   - p: The running Bubble Tea program, used to send streaming messages.
//   - moduleName: The module being uninstalled.
//   - commands: Shell commands to reverse the installation.
//   - dotfilesDir: The dotfiles repository root (for unstow).
//   - stowEnabled: Whether to remove stow symlinks after commands run.
func RunUninstallStreaming(p *tea.Program, moduleName string, commands []string, dotfilesDir string, stowEnabled bool) tea.Cmd {
	return func() tea.Msg {
		p.Send(UninstallStartMsg{ModuleName: moduleName})

		// Run each uninstall command sequentially.
		for i, cmdStr := range commands {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       fmt.Sprintf("\n▸ Uninstall step %d/%d: %s", i+1, len(commands), cmdStr),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			err := utils.RunCommandStreaming(ctx, cmdStr, func(line string, isStderr bool) {
				p.Send(InstallOutputMsg{
					ModuleName: moduleName,
					Line:       line,
					IsStderr:   isStderr,
				})
			})
			cancel()

			if err != nil {
				return UninstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("step %d/%d %q: %w", i+1, len(commands), cmdStr, err),
				}
			}
		}

		// Remove stow symlinks after all uninstall commands succeed.
		if stowEnabled {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "\n▸ Removing stow symlinks...",
			})
			if err := utils.Unstow(moduleName, dotfilesDir); err != nil {
				return UninstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("unstow: %w", err),
				}
			}
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "✓ Stow links removed",
			})
		}

		// Update the tracking file to mark this module as uninstalled.
		if err := utils.SetModuleUninstalled(moduleName); err != nil {
			return UninstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("tracking uninstall: %w", err),
			}
		}

		return UninstallCompleteMsg{
			ModuleName: moduleName,
			Success:    true,
		}
	}
}
