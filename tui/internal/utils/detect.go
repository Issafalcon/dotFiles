// This file provides software detection utilities for the TUI application.
//
// These helpers are used to check if prerequisite tools are installed,
// get their versions, and detect the runtime environment (e.g., WSL).
//
// Key Go concepts used here:
//   - exec.LookPath: Searching PATH for executables (https://pkg.go.dev/os/exec#LookPath)
//   - os.ReadFile: Reading entire file contents (https://pkg.go.dev/os#ReadFile)
//   - strings: String manipulation functions (https://pkg.go.dev/strings)
package utils

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// IsCommandAvailable checks if a command is available in the system's PATH.
//
// It uses exec.LookPath, which searches through the directories in the PATH
// environment variable for an executable matching the given name.
//
// Parameters:
//   - command: The command name to look for (e.g., "git", "stow", "nvim").
//
// Returns:
//   - bool: true if the command is found in PATH, false otherwise.
//
// In Go, bool is a built-in type with values true and false.
// See: https://go.dev/ref/spec#Boolean_types
// See: https://pkg.go.dev/os/exec#LookPath
func IsCommandAvailable(command string) bool {
	// exec.LookPath returns the full path to the executable and an error.
	// We only care about whether it succeeded (err == nil), so we discard
	// the path using the blank identifier (_).
	// See: https://go.dev/doc/effective_go#blank
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandVersion runs `command --version` and returns the first line of output.
//
// Most CLI tools print their version information when invoked with --version.
// We capture the first line, which typically contains the version string.
//
// Parameters:
//   - command: The command to query for its version (e.g., "git", "node").
//
// Returns:
//   - string: The first line of the version output.
//   - error: An error if the command fails or produces no output.
//
// Example:
//
//	version, err := GetCommandVersion("git")
//	// version might be: "git version 2.43.0"
func GetCommandVersion(command string) (string, error) {
	// Use a short timeout — version checks should be nearly instantaneous.
	// context.WithTimeout creates a context that automatically cancels after
	// the given duration. The cancel function must be called to release resources.
	// See: https://pkg.go.dev/context#WithTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use our RunCommand helper to execute the version check.
	// The command string is passed through a shell, so it handles PATH lookup.
	result, err := RunCommand(ctx, command+" --version")
	if err != nil {
		return "", err
	}

	// Check if we got any stdout output.
	// len() is a built-in function that returns the length of slices, strings,
	// maps, and channels.
	// See: https://pkg.go.dev/builtin#len
	if len(result.Stdout) > 0 {
		return result.Stdout[0], nil
	}

	// Some commands print version info to stderr instead of stdout.
	if len(result.Stderr) > 0 {
		return result.Stderr[0], nil
	}

	return "", nil
}

// DetectWSL checks if the application is running inside Windows Subsystem for Linux.
//
// WSL detection is done by reading /proc/version, which on WSL contains
// "Microsoft" or "WSL" in the kernel version string.
//
// On a native Linux system, /proc/version looks like:
//
//	Linux version 6.1.0-18-amd64 (debian-kernel@lists.debian.org) ...
//
// On WSL, it looks like:
//
//	Linux version 5.15.133.1-microsoft-standard-WSL2 ...
//
// Returns:
//   - bool: true if running in WSL, false otherwise.
func DetectWSL() bool {
	// os.ReadFile reads an entire file into memory as a byte slice ([]byte).
	// This is convenient for small files like /proc/version.
	// See: https://pkg.go.dev/os#ReadFile
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		// If /proc/version doesn't exist (e.g., on macOS), we're not in WSL.
		return false
	}

	// Convert []byte to string for text operations.
	// In Go, string() is a type conversion, not a function call.
	// Strings and byte slices share the same underlying data representation (UTF-8).
	// See: https://go.dev/blog/strings
	// See: https://go.dev/ref/spec#Conversions
	content := string(data)

	// strings.Contains checks if a string contains a substring.
	// See: https://pkg.go.dev/strings#Contains
	//
	// strings.ToLower converts to lowercase for case-insensitive matching.
	// See: https://pkg.go.dev/strings#ToLower
	lowered := strings.ToLower(content)
	return strings.Contains(lowered, "microsoft") || strings.Contains(lowered, "wsl")
}
