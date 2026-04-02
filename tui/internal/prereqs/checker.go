// Package prereqs handles checking and installing system prerequisites.
//
// Before the TUI can manage dotfiles, certain tools must be installed on the
// system (git, stow, curl, etc.). This package defines what those prerequisites
// are, how to check if they're present, and how to install them.
//
// This file (checker.go) contains the data definitions and pure logic.
// The companion file (prereqs.go) contains the Bubble Tea model for the UI.
//
// # Key Go Concepts Used
//
//   - Structs: Custom data types that group related fields
//     See: https://go.dev/tour/moretypes/2
//   - Slices: Dynamic arrays — the most common collection type in Go
//     See: https://go.dev/doc/effective_go#slices
//   - os/exec: Running external commands from Go
//     See: https://pkg.go.dev/os/exec
//
// For more on Go basics: https://go.dev/doc/
package prereqs

import (
	"os/exec"
	"strings"
)

// Prereq represents a single system prerequisite that the dotfiles manager needs.
//
// In Go, a struct is a composite type that groups named fields together.
// Each field has a name and a type. Exported fields (capitalised) can be
// accessed from other packages; unexported fields (lowercase) are private.
//
// Struct fields are accessed with dot notation: p.Name, p.Description, etc.
//
// See: https://go.dev/ref/spec#Struct_types
// See: https://go.dev/tour/moretypes/2
type Prereq struct {
	// Name is the binary or package name (e.g., "git", "curl").
	Name string

	// Description is a human-readable summary of what this tool does.
	Description string

	// CheckCommand is the shell command used to verify the tool is installed.
	// For example, "command -v git" checks if the 'git' binary is on the PATH.
	// The special shell builtin 'command -v' is more portable than 'which'.
	// See: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/command.html
	CheckCommand string

	// DocsURL is a link to the official documentation for this tool.
	DocsURL string

	// InstallCommand is the apt package name used to install this tool.
	// For example, "git" becomes "sudo apt install git".
	InstallCommand string
}

// GetRequiredPrereqs returns the full list of system prerequisites.
//
// This list mirrors the packages installed by prerequisites.sh in the
// dotfiles repository root. Each entry specifies how to check for the
// tool and how to install it if missing.
//
// In Go, functions that return data (no receiver) are called "package-level
// functions". They're similar to static methods in other languages.
// The []Prereq return type is a "slice of Prereq" — Go's dynamic array.
//
// See: https://go.dev/tour/moretypes/7 (slices)
// See: https://go.dev/doc/effective_go#composite_literals (struct literals)
func GetRequiredPrereqs() []Prereq {
	// A slice literal creates and initialises a slice in one expression.
	// Each {} block is a struct literal — you can name the fields explicitly
	// (Name: "git") for clarity, which is preferred over positional arguments.
	// See: https://go.dev/ref/spec#Composite_literals
	return []Prereq{
		{
			Name:           "git",
			Description:    "Distributed version control system",
			CheckCommand:   "command -v git",
			DocsURL:        "https://git-scm.com/doc",
			InstallCommand: "git",
		},
		{
			Name:           "stow",
			Description:    "Symlink farm manager for dotfiles",
			CheckCommand:   "command -v stow",
			DocsURL:        "https://www.gnu.org/software/stow/",
			InstallCommand: "stow",
		},
		{
			Name:           "zsh",
			Description:    "Z shell — powerful interactive shell",
			CheckCommand:   "command -v zsh",
			DocsURL:        "https://www.zsh.org/",
			InstallCommand: "zsh",
		},
		{
			Name:           "curl",
			Description:    "Command-line URL transfer tool",
			CheckCommand:   "command -v curl",
			DocsURL:        "https://curl.se/docs/",
			InstallCommand: "curl",
		},
		{
			Name:           "wget",
			Description:    "Non-interactive network file downloader",
			CheckCommand:   "command -v wget",
			DocsURL:        "https://www.gnu.org/software/wget/",
			InstallCommand: "wget",
		},
		{
			Name:           "zip",
			Description:    "Compression and archive utility",
			CheckCommand:   "command -v zip",
			DocsURL:        "https://infozip.sourceforge.net/",
			InstallCommand: "zip",
		},
		{
			Name:           "unzip",
			Description:    "Extraction utility for zip archives",
			CheckCommand:   "command -v unzip",
			DocsURL:        "https://infozip.sourceforge.net/",
			InstallCommand: "unzip",
		},
		{
			Name:           "build-essential",
			Description:    "C/C++ compiler and build tools (gcc, make)",
			CheckCommand:   "dpkg -s build-essential",
			DocsURL:        "https://packages.debian.org/build-essential",
			InstallCommand: "build-essential",
		},
		{
			Name:           "libssl-dev",
			Description:    "SSL/TLS development headers and libraries",
			CheckCommand:   "dpkg -s libssl-dev",
			DocsURL:        "https://www.openssl.org/docs/",
			InstallCommand: "libssl-dev",
		},
		{
			Name:           "jq",
			Description:    "Lightweight command-line JSON processor",
			CheckCommand:   "command -v jq",
			DocsURL:        "https://jqlang.github.io/jq/",
			InstallCommand: "jq",
		},
		{
			Name:           "fdclone",
			Description:    "Console-based file manager",
			CheckCommand:   "dpkg -s fdclone",
			DocsURL:        "https://hp.vector.co.jp/authors/VA012337/soft/fd/",
			InstallCommand: "fdclone",
		},
	}
}

// CheckPrereq runs the check command for a single prerequisite and returns
// true if the tool is installed, false otherwise.
//
// # How exec.Command Works
//
// exec.Command creates a new *exec.Cmd struct that represents an external
// command to be run. We use "sh -c <command>" to run the check through a
// shell, because some checks use shell builtins like "command -v" which
// aren't standalone binaries.
//
//   - exec.Command("sh", "-c", "command -v git") → runs: sh -c "command -v git"
//   - cmd.Run() executes the command and waits for it to finish
//   - Run() returns nil on success (exit code 0) or an error on failure
//
// See: https://pkg.go.dev/os/exec#Command
// See: https://pkg.go.dev/os/exec#Cmd.Run
func CheckPrereq(p Prereq) bool {
	// exec.Command returns an *exec.Cmd — a pointer to a Cmd struct.
	// In Go, the & operator takes the address of a value, and * dereferences it.
	// The exec package returns a pointer so the caller can configure the command
	// (e.g., set Stdin, Stdout, Env) before running it.
	// See: https://go.dev/tour/moretypes/1 (pointers)
	cmd := exec.Command("sh", "-c", p.CheckCommand)

	// cmd.Run() both starts the command and waits for completion.
	// It returns an *exec.ExitError if the command exits with a non-zero code,
	// or another error type if the command couldn't be started at all.
	// We only care whether it succeeded (err == nil) or not.
	err := cmd.Run()

	// In Go, nil is the zero value for pointers, interfaces, maps, slices,
	// channels, and function types. A nil error means "no error" / success.
	// See: https://go.dev/doc/effective_go#errors
	return err == nil
}

// InstallPrereqs generates the shell commands needed to install all missing
// prerequisites. It returns a slice of command strings that can be executed
// in sequence.
//
// This function doesn't execute the commands directly — it just builds them.
// The caller (the Bubble Tea model) decides when and how to run them.
//
// # Why Return Commands Instead of Running Them?
//
// Running "sudo apt install" requires elevated privileges and may prompt for
// a password. By returning the command strings, the UI layer can decide how
// to present them (e.g., show them to the user, run them via tea.Exec, etc.).
//
// See: https://go.dev/doc/effective_go#slices (working with slices)
func InstallPrereqs(missing []Prereq) []string {
	// len() is a built-in function that returns the length of various types:
	// slices, maps, strings, arrays, and channels.
	// See: https://pkg.go.dev/builtin#len
	if len(missing) == 0 {
		// Return nil for an empty slice. In Go, a nil slice and an empty slice
		// behave the same for most operations (len, range, append all work).
		// See: https://go.dev/doc/effective_go#slices
		return nil
	}

	// make() is a built-in that creates slices, maps, and channels.
	// make([]string, 0, len(missing)) creates a string slice with:
	//   - length 0 (no elements yet)
	//   - capacity len(missing) (pre-allocated space for efficiency)
	// This avoids repeated memory allocations as we append.
	// See: https://pkg.go.dev/builtin#make
	// See: https://go.dev/blog/slices-intro
	packages := make([]string, 0, len(missing))

	// range iterates over slices, maps, strings, and channels.
	// For slices, it returns (index, value) on each iteration.
	// The _ (blank identifier) discards the index since we don't need it.
	// See: https://go.dev/tour/moretypes/16
	for _, p := range missing {
		packages = append(packages, p.InstallCommand)
	}

	// strings.Join concatenates slice elements with a separator.
	// e.g., strings.Join(["git", "curl", "wget"], " ") → "git curl wget"
	// See: https://pkg.go.dev/strings#Join
	return []string{
		"sudo apt update",
		"sudo apt install -y " + strings.Join(packages, " "),
	}
}
