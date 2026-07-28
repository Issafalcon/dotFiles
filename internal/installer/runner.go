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
	"os"
	"os/exec"
	"path/filepath"
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

// NeedsSudoScript returns true if the script at path contains "sudo".
// Missing or unreadable scripts are treated as not needing sudo.
func NeedsSudoScript(scriptPath string) bool {
	if scriptPath == "" {
		return false
	}
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "sudo")
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

// RunInstallInteractive suspends the TUI and runs install.sh with a real TTY
// (stdin/stdout/stderr). Use for modules with requires_input: true so prompts
// (rustup, installers, etc.) work. After the process exits, stow + tracking run
// and InstallCompleteMsg is returned so the dep queue can continue.
func RunInstallInteractive(moduleName, scriptPath, modulesDir string, stowEnabled bool) tea.Cmd {
	c := exec.Command("bash", scriptPath)
	c.Env = utils.BrewAwareEnv()
	c.Dir = filepath.Dir(scriptPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("install.sh: %w", err),
			}
		}
		if stowEnabled {
			if err := utils.Stow(moduleName, modulesDir); err != nil {
				return InstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("stow: %w", err),
				}
			}
		}
		if err := utils.SetModuleInstalled(moduleName); err != nil {
			return InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("tracking install: %w", err),
			}
		}
		return InstallCompleteMsg{ModuleName: moduleName, Success: true}
	})
}

// scriptCommand builds a bash invocation for a module install/uninstall script.
func scriptCommand(scriptPath string) string {
	return "bash " + shellQuote(scriptPath)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// --- Streaming install ---

// RunInstallStreaming returns a tea.Cmd that runs the module's install.sh
// (if present) in a background goroutine, streaming each line of output to
// the UI via p.Send(). The final message returned by the tea.Cmd is
// InstallCompleteMsg.
//
// Parameters:
//   - p: The running Bubble Tea program, used to send streaming messages.
//   - moduleName: The module being installed.
//   - scriptPath: Absolute path to install.sh, or "" for stow-only modules.
//   - modulesDir: The modules/ directory (stow directory).
//   - stowEnabled: Whether to run stow after the script succeeds.
func RunInstallStreaming(p *tea.Program, moduleName string, scriptPath string, modulesDir string, stowEnabled bool) tea.Cmd {
	return func() tea.Msg {
		if p == nil {
			return InstallCompleteMsg{
				ModuleName: moduleName,
				Success:    false,
				Error:      fmt.Errorf("internal: tea program not ready"),
			}
		}
		p.Send(InstallStartMsg{ModuleName: moduleName})

		if scriptPath != "" {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       fmt.Sprintf("\n▸ Running %s", scriptPath),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			// Refresh sudo timestamp while long scripts run (default timeout ~15m).
			// Without this, a late `sudo` in install.sh hangs with no TTY after brew.
			if NeedsSudoScript(scriptPath) {
				stopKeepAlive := keepSudoAlive(ctx)
				defer stopKeepAlive()
			}

			err := utils.RunCommandStreaming(ctx, scriptCommand(scriptPath), func(line string, isStderr bool) {
				p.Send(InstallOutputMsg{
					ModuleName: moduleName,
					Line:       line,
					IsStderr:   isStderr,
				})
			})

			if err != nil {
				return InstallCompleteMsg{
					ModuleName: moduleName,
					Success:    false,
					Error:      fmt.Errorf("install.sh: %w", err),
				}
			}
		}

		if stowEnabled {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "\n▸ Creating stow symlinks...",
			})
			if err := utils.Stow(moduleName, modulesDir); err != nil {
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

// keepSudoAlive periodically runs `sudo -n true` so an earlier sudo -v stays
// valid across long install scripts. Returns a stop function.
func keepSudoAlive(ctx context.Context) func() {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-ticker.C:
				_ = exec.CommandContext(ctx, "sudo", "-n", "true").Run()
			}
		}
	}()
	return func() { close(done) }
}

// RunInstallWithSend executes the install script and sends progress messages via
// a provided send function. This is used by the Orchestrator for real-time
// streaming output to the Bubble Tea UI.
func RunInstallWithSend(ctx context.Context, moduleName string, scriptPath string, modulesDir string, stowEnabled bool, send func(tea.Msg)) {
	send(InstallStartMsg{ModuleName: moduleName})

	if scriptPath != "" {
		send(InstallProgressMsg{
			ModuleName: moduleName,
			Step:       1,
			TotalSteps: 1,
		})

		err := utils.RunCommandStreaming(ctx, scriptCommand(scriptPath), func(line string, isStderr bool) {
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
				Error:      fmt.Errorf("install.sh failed: %w", err),
			})
			return
		}
	}

	if stowEnabled {
		if err := utils.Stow(moduleName, modulesDir); err != nil {
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

// RunUninstallStreaming returns a tea.Cmd that runs uninstall.sh (if present),
// then removes stow symlinks and updates the tracking file.
func RunUninstallStreaming(p *tea.Program, moduleName string, scriptPath string, modulesDir string, stowEnabled bool) tea.Cmd {
	return func() tea.Msg {
		p.Send(UninstallStartMsg{ModuleName: moduleName})

		if scriptPath != "" {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       fmt.Sprintf("\n▸ Running %s", scriptPath),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			err := utils.RunCommandStreaming(ctx, scriptCommand(scriptPath), func(line string, isStderr bool) {
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
					Error:      fmt.Errorf("uninstall.sh: %w", err),
				}
			}
		}

		if stowEnabled {
			p.Send(InstallOutputMsg{
				ModuleName: moduleName,
				Line:       "\n▸ Removing stow symlinks...",
			})
			if err := utils.Unstow(moduleName, modulesDir); err != nil {
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
