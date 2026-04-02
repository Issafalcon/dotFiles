// This file provides a wrapper around GNU Stow for managing dotfile symlinks,
// and functions for tracking which modules are installed.
//
// GNU Stow is a symlink farm manager. It creates symlinks from a "stow directory"
// (where your dotfiles live) to a "target directory" (usually your home directory).
// See: https://www.gnu.org/software/stow/
//
// Key Go concepts used here:
//   - os package: File system operations (https://pkg.go.dev/os)
//   - filepath: Platform-independent path manipulation (https://pkg.go.dev/path/filepath)
//   - Error wrapping with %w: Preserves the original error for inspection (https://pkg.go.dev/fmt#Errorf)
//   - Variadic error handling: Go's explicit error return pattern
package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// dotFileModulesFile is the name of the file that tracks installed modules.
// It is stored in the user's home directory.
//
// In Go, constants are declared with the const keyword. They must be compile-time
// evaluable (no function calls allowed). String constants are untyped by default.
// See: https://go.dev/ref/spec#Constants
// See: https://go.dev/doc/effective_go#constants
const dotFileModulesFile = ".dotFileModules"

// Stow runs GNU Stow to create symlinks for the given module.
//
// It executes: stow <moduleName> from within the dotfilesDir.
// The --verbose flag provides output for debugging.
// The --target flag explicitly sets the home directory as the target.
//
// Parameters:
//   - moduleName: The name of the module directory to stow (e.g., "nvim", "zsh").
//   - dotfilesDir: The path to the dotfiles repository root.
//
// Returns:
//   - error: nil on success, or an error describing what went wrong.
//     In Go, returning error as the last value is a strong convention.
//     See: https://go.dev/doc/effective_go#errors
func Stow(moduleName string, dotfilesDir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// fmt.Errorf with %w wraps the original error, preserving it for
		// errors.Is() and errors.Unwrap() inspection by callers.
		// See: https://pkg.go.dev/fmt#Errorf
		// See: https://go.dev/blog/go1.13-errors
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Build the stow command with explicit target directory.
	// fmt.Sprintf formats a string using printf-style verbs.
	// See: https://pkg.go.dev/fmt#Sprintf
	command := fmt.Sprintf("cd %s && stow --verbose --target=%s %s",
		shellEscape(dotfilesDir),
		shellEscape(homeDir),
		shellEscape(moduleName),
	)

	// Use a 60-second timeout for stow operations.
	// context.WithTimeout returns a derived context that automatically cancels
	// after the specified duration, plus a cancel function you must call to
	// release resources (even if the timeout hasn't expired).
	// See: https://pkg.go.dev/context#WithTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	// defer ensures cancel() runs when Stow returns, preventing context leaks.
	// See: https://go.dev/doc/effective_go#defer
	defer cancel()

	result, err := RunCommand(ctx, command)
	if err != nil {
		return fmt.Errorf("running stow for module %q: %w", moduleName, err)
	}

	if result.ExitCode != 0 {
		// strings.Join concatenates slice elements with a separator.
		// See: https://pkg.go.dev/strings#Join
		return fmt.Errorf("stow failed for module %q (exit code %d): %s",
			moduleName, result.ExitCode, strings.Join(result.Stderr, "\n"))
	}

	return nil
}

// Unstow runs GNU Stow with the -D flag to remove symlinks for the given module.
//
// It executes: stow -D <moduleName> from within the dotfilesDir.
func Unstow(moduleName string, dotfilesDir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// The -D flag tells stow to "unstow" (delete symlinks) instead of creating them.
	command := fmt.Sprintf("cd %s && stow --verbose -D --target=%s %s",
		shellEscape(dotfilesDir),
		shellEscape(homeDir),
		shellEscape(moduleName),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := RunCommand(ctx, command)
	if err != nil {
		return fmt.Errorf("running unstow for module %q: %w", moduleName, err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("unstow failed for module %q (exit code %d): %s",
			moduleName, result.ExitCode, strings.Join(result.Stderr, "\n"))
	}

	return nil
}

// GetInstalledModules reads the ~/.dotFileModules file and returns the list of
// currently installed module names.
//
// Each line in the file is one module name. Empty lines and whitespace are trimmed.
//
// Returns:
//   - []string: A slice of installed module names.
//   - error: An error if the file cannot be read (returns empty slice if file doesn't exist).
func GetInstalledModules() ([]string, error) {
	modulesPath, err := getModulesFilePath()
	if err != nil {
		return nil, err
	}

	// os.Open opens a file for reading. It returns an *os.File and an error.
	// See: https://pkg.go.dev/os#Open
	file, err := os.Open(modulesPath)
	if err != nil {
		// os.IsNotExist checks if an error indicates the file doesn't exist.
		// If the tracking file doesn't exist yet, no modules are installed.
		// See: https://pkg.go.dev/os#IsNotExist
		if os.IsNotExist(err) {
			// Return nil slice (zero value for slices) and no error.
			// nil slices behave like empty slices for most operations.
			// See: https://go.dev/doc/effective_go#slices
			return nil, nil
		}
		return nil, fmt.Errorf("opening modules file: %w", err)
	}

	// defer file.Close() ensures the file is closed when the function returns,
	// regardless of how it returns (normal return, early error return, etc.).
	// This is the idiomatic Go pattern for resource cleanup.
	// See: https://go.dev/doc/effective_go#defer
	defer file.Close()

	var modules []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// strings.TrimSpace removes leading and trailing whitespace.
		// See: https://pkg.go.dev/strings#TrimSpace
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			modules = append(modules, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading modules file: %w", err)
	}

	return modules, nil
}

// SetModuleInstalled appends a module name to the ~/.dotFileModules file if it
// is not already present. This is idempotent — calling it twice with the same
// name has no additional effect.
func SetModuleInstalled(name string) error {
	// First, check if the module is already tracked.
	modules, err := GetInstalledModules()
	if err != nil {
		return fmt.Errorf("checking installed modules: %w", err)
	}

	// range iterates over slices, maps, strings, and channels.
	// For slices, it yields (index, value) pairs. The blank identifier (_)
	// discards the index since we only need the value.
	// See: https://go.dev/ref/spec#For_range
	// See: https://go.dev/doc/effective_go#for
	for _, m := range modules {
		if m == name {
			// Already installed, nothing to do. Return nil (no error).
			return nil
		}
	}

	modulesPath, err := getModulesFilePath()
	if err != nil {
		return err
	}

	// os.OpenFile gives fine-grained control over how a file is opened.
	// os.O_APPEND: Writes go to end of file.
	// os.O_CREATE: Create the file if it doesn't exist.
	// os.O_WRONLY: Open for writing only.
	// 0644: Unix file permissions (owner read/write, group/others read-only).
	// See: https://pkg.go.dev/os#OpenFile
	// See: https://pkg.go.dev/os#pkg-constants (for O_APPEND, O_CREATE, etc.)
	file, err := os.OpenFile(modulesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening modules file for writing: %w", err)
	}
	defer file.Close()

	// fmt.Fprintln writes a line (with newline) to the given io.Writer.
	// *os.File implements io.Writer, so it can be passed directly.
	// See: https://pkg.go.dev/fmt#Fprintln
	if _, err := fmt.Fprintln(file, name); err != nil {
		return fmt.Errorf("writing module name: %w", err)
	}

	return nil
}

// SetModuleUninstalled removes a module name from the ~/.dotFileModules file.
// It rewrites the file without the specified module name.
func SetModuleUninstalled(name string) error {
	modules, err := GetInstalledModules()
	if err != nil {
		return fmt.Errorf("reading installed modules: %w", err)
	}

	// Build a new list excluding the module to remove.
	// We pre-allocate with make() to avoid repeated allocations during append.
	//
	// make([]T, length, capacity) creates a slice with the given length and capacity.
	// Here length=0 (starts empty) and capacity=len(modules) (pre-allocates memory).
	// See: https://go.dev/doc/effective_go#allocation_make
	// See: https://pkg.go.dev/builtin#make
	filtered := make([]string, 0, len(modules))
	for _, m := range modules {
		if m != name {
			// append adds elements to a slice, growing it as needed.
			// It may allocate a new underlying array if capacity is exceeded.
			// See: https://pkg.go.dev/builtin#append
			filtered = append(filtered, m)
		}
	}

	modulesPath, err := getModulesFilePath()
	if err != nil {
		return err
	}

	// Join all module names with newlines and add a trailing newline.
	content := strings.Join(filtered, "\n")
	if len(filtered) > 0 {
		content += "\n"
	}

	// os.WriteFile atomically writes data to a file (or creates it).
	// []byte(content) converts the string to a byte slice.
	// In Go, strings are immutable UTF-8 byte sequences, and []byte is mutable.
	// See: https://pkg.go.dev/os#WriteFile
	// See: https://go.dev/blog/strings
	if err := os.WriteFile(modulesPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing modules file: %w", err)
	}

	return nil
}

// GetDotfilesDir returns the path to the dotfiles repository root.
//
// Detection order:
//  1. DOTFILES_DIR environment variable (if set).
//  2. The parent directory of the running executable's location.
//     This works because the TUI binary lives inside the tui/ subdirectory
//     of the dotfiles repo.
//  3. Falls back to ~/dotFiles as a sensible default.
func GetDotfilesDir() string {
	// os.Getenv reads an environment variable. It returns "" if not set.
	// See: https://pkg.go.dev/os#Getenv
	if envDir := os.Getenv("DOTFILES_DIR"); envDir != "" {
		return envDir
	}

	// os.Executable returns the path of the running executable.
	// See: https://pkg.go.dev/os#Executable
	execPath, err := os.Executable()
	if err == nil {
		// filepath.Dir returns the directory portion of a path.
		// We go up two levels: from tui/binary -> tui/ -> dotfiles root.
		// See: https://pkg.go.dev/path/filepath#Dir
		tuiDir := filepath.Dir(execPath)
		dotfilesDir := filepath.Dir(tuiDir)

		// os.Stat returns file info. We check if the directory exists and looks
		// like a dotfiles repo (contains a bootstrap.sh file).
		if _, statErr := os.Stat(filepath.Join(dotfilesDir, "bootstrap.sh")); statErr == nil {
			return dotfilesDir
		}
	}

	// Fall back to ~/dotFiles.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't even get the home directory, return a relative path.
		return "."
	}
	return filepath.Join(homeDir, "dotFiles")
}

// getModulesFilePath returns the full path to the ~/.dotFileModules tracking file.
func getModulesFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, dotFileModulesFile), nil
}

// shellEscape wraps a path in single quotes to prevent shell interpretation
// of special characters like spaces and glob patterns.
//
// This is a simple approach — for production code, consider using exec.Command
// with separate arguments instead of shell string concatenation.
func shellEscape(s string) string {
	// strings.ReplaceAll replaces all occurrences of a substring.
	// We escape single quotes within the string by ending the quote,
	// adding an escaped quote, and re-opening the quote.
	// See: https://pkg.go.dev/strings#ReplaceAll
	escaped := strings.ReplaceAll(s, "'", "'\\''")
	return "'" + escaped + "'"
}
