// This file provides cross-platform URL opening for the TUI application.
//
// On Linux, URLs are opened with xdg-open. In WSL (Windows Subsystem for Linux),
// we detect the environment and use appropriate Windows-side commands.
//
// Key Go concepts used here:
//   - runtime.GOOS: Compile-time OS detection (https://pkg.go.dev/runtime#pkg-constants)
//   - exec.LookPath: Finding executables in PATH (https://pkg.go.dev/os/exec#LookPath)
//   - exec.Command: Running external commands (https://pkg.go.dev/os/exec#Command)
package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL opens the given URL in the user's default browser.
//
// It detects the platform and WSL environment to choose the right command:
//   - Linux (native): uses xdg-open (https://linux.die.net/man/1/xdg-open)
//   - WSL: uses wslview (from wslu package) or falls back to cmd.exe /c start
//   - macOS: uses open (included for completeness, though this TUI targets Linux)
//
// Parameters:
//   - url: The URL string to open (e.g., "https://github.com/issafalcon/dotFiles").
//
// Returns:
//   - error: nil on success, or an error if the URL could not be opened.
func OpenURL(url string) error {
	// First, check if we're running inside WSL (Windows Subsystem for Linux).
	// WSL needs special handling because Linux browser openers won't work —
	// we need to call into the Windows side to open the browser.
	if DetectWSL() {
		return openURLInWSL(url)
	}

	// runtime.GOOS is a compile-time constant that identifies the operating system.
	// Possible values: "linux", "darwin", "windows", "freebsd", etc.
	// See: https://pkg.go.dev/runtime#pkg-constants
	// See: https://go.dev/doc/install/source#environment (list of GOOS values)
	switch runtime.GOOS {
	case "linux":
		return openWithCommand("xdg-open", url)
	case "darwin":
		return openWithCommand("open", url)
	default:
		// %q formats a string with Go-syntax quoting (adds double quotes and escapes).
		return fmt.Errorf("unsupported platform: %q", runtime.GOOS)
	}
}

// openURLInWSL tries multiple strategies to open a URL from within WSL.
//
// WSL Strategy:
//  1. Try wslview (from the wslu utilities package) — this is the recommended way.
//  2. Fall back to cmd.exe /c start — directly invokes Windows command processor.
func openURLInWSL(url string) error {
	// exec.LookPath searches for an executable in the directories listed in
	// the PATH environment variable. It returns the full path if found.
	// See: https://pkg.go.dev/os/exec#LookPath
	if _, err := exec.LookPath("wslview"); err == nil {
		return openWithCommand("wslview", url)
	}

	// Fall back to cmd.exe. In WSL, Windows executables (like cmd.exe) are
	// accessible through the PATH. The /c flag tells cmd.exe to execute
	// the following command and then terminate.
	if _, err := exec.LookPath("cmd.exe"); err == nil {
		// cmd.exe /c start "" "url" — the empty string "" is needed as the
		// window title parameter for the start command when the URL contains
		// special characters.
		cmd := exec.Command("cmd.exe", "/c", "start", "", url)
		return cmd.Run()
	}

	return fmt.Errorf("no suitable browser opener found in WSL (tried wslview and cmd.exe)")
}

// openWithCommand runs the specified command with the URL as an argument.
//
// exec.Command creates a Cmd struct representing an external command.
// Unlike exec.CommandContext, this version has no context for cancellation.
// cmd.Run() starts the command and waits for it to complete.
// See: https://pkg.go.dev/os/exec#Command
// See: https://pkg.go.dev/os/exec#Cmd.Run
func openWithCommand(command string, url string) error {
	cmd := exec.Command(command, url)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("opening URL with %s: %w", command, err)
	}

	return nil
}
