// Package module defines the core data structures for dotfile modules.
//
// A "module" in this context is a single tool or application whose configuration
// is managed by this dotfiles repository. Each module has metadata (name, icon,
// description), installation commands (from the original install.sh scripts),
// and dependency information for ordering installations.
//
// # Go Structs
//
// Go uses structs instead of classes. A struct is a collection of typed fields.
// Unlike OOP languages, Go structs don't have constructors or inheritance.
// Instead, you compose structs together and attach methods via receiver functions.
//
// See: https://go.dev/doc/effective_go#composite_literals
// See: https://go.dev/tour/moretypes/2
// See: https://go.dev/ref/spec#Struct_types
//
// # Custom Types and Enums
//
// Go doesn't have a built-in enum keyword. Instead, you create a new type based
// on an underlying type (usually int or string) and define constants using iota.
//
// See: https://go.dev/ref/spec#Iota
// See: https://go.dev/ref/spec#Constant_declarations
package module

// InstallStatus represents the current installation state of a module.
// This is a custom type based on int — Go's idiomatic way to create enumerations.
//
// The type keyword creates a new named type. Even though InstallStatus is based
// on int, Go's type system treats it as a distinct type — you can't accidentally
// assign a plain int to it without an explicit conversion.
//
// See: https://go.dev/ref/spec#Type_definitions
// See: https://go.dev/doc/effective_go#constants
type InstallStatus int

// These constants define all possible installation states using iota.
//
// iota is a special Go constant generator. Within a const block, iota starts
// at 0 and increments by 1 for each constant. This gives us:
//   - StatusUnknown     = 0
//   - StatusInstalled   = 1
//   - StatusNotInstalled = 2
//   - StatusInstalling  = 3
//   - StatusFailed      = 4
//
// See: https://go.dev/ref/spec#Iota
const (
	StatusUnknown      InstallStatus = iota // Installation status has not been checked yet.
	StatusInstalled                         // Module is confirmed installed on the system.
	StatusNotInstalled                      // Module is confirmed NOT installed.
	StatusInstalling                        // Module installation is currently in progress.
	StatusFailed                            // Module installation was attempted but failed.
)

// String returns a human-readable label for an InstallStatus value.
//
// This method satisfies the fmt.Stringer interface, which means any time you
// pass an InstallStatus to fmt.Println, fmt.Sprintf, etc., Go will
// automatically call this method to get the string representation.
//
// The (s InstallStatus) part is called a "receiver" — it attaches this function
// to the InstallStatus type, making it a method rather than a standalone function.
//
// See: https://go.dev/tour/methods/1
// See: https://pkg.go.dev/fmt#Stringer
func (s InstallStatus) String() string {
	// A switch statement in Go doesn't need "break" — each case automatically
	// breaks unless you use the "fallthrough" keyword.
	// See: https://go.dev/tour/flowcontrol/9
	switch s {
	case StatusUnknown:
		return "Unknown"
	case StatusInstalled:
		return "Installed"
	case StatusNotInstalled:
		return "Not Installed"
	case StatusInstalling:
		return "Installing..."
	case StatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// ExternalDep represents an external tool or binary that a module requires
// but which is NOT another module in this dotfiles repo.
//
// For example, nvim might need "ripgrep" installed, and we need to know how
// to check for it and install it if missing.
//
// Each field has a specific purpose:
//   - Name: Human-readable name of the dependency (e.g., "ripgrep")
//   - CheckCommand: Shell command to verify it's installed (e.g., "rg --version")
//   - InstallCommand: Shell command to install it (e.g., "sudo apt install ripgrep")
//   - InstallMethod: Which package manager to use (apt, brew, cargo, npm, pip, curl)
//
// See: https://go.dev/ref/spec#Struct_types
type ExternalDep struct {
	Name           string // Human-readable name of the external tool.
	CheckCommand   string // Shell command to verify the tool is installed.
	InstallCommand string // Shell command to install the tool if missing.
	InstallMethod  string // Package manager used: "apt", "brew", "cargo", "npm", "pip", or "curl".
}

// ConfigOption represents a user-configurable choice for a module.
//
// Some modules need user input during setup (e.g., "which .NET version to install?").
// ConfigOption defines what the choice is, what the default value is, and what
// values are valid.
//
// The Choices slice can be nil/empty if the option accepts freeform text input.
//
// # Slices in Go
//
// []string is a "slice" — Go's dynamically-sized array type. Slices are
// reference types backed by an underlying array. A nil slice ([]string(nil))
// and an empty slice ([]string{}) both have length 0 but behave slightly
// differently with JSON serialization.
//
// See: https://go.dev/tour/moretypes/7
// See: https://go.dev/blog/slices-intro
type ConfigOption struct {
	Name        string   // Short identifier for this option (e.g., "dotnet_version").
	Description string   // Human-readable description shown to the user.
	Default     string   // Default value if the user doesn't choose.
	Choices     []string // Valid values the user can pick from. Empty means freeform.
}

// Module represents a single dotfile module — one tool, application, or
// configuration managed by this repository.
//
// This is the central data structure of the entire application. Each Module
// carries everything needed to display it in the TUI, check its status,
// install it, and uninstall it.
//
// # Struct Field Ordering
//
// Fields are ordered by conceptual grouping: identity fields first, then
// categorization, then installation details, then runtime state. This makes
// the struct easier to read and reason about.
//
// # Pointer vs Value Receivers
//
// We generally pass *Module (pointer to Module) rather than Module (value copy)
// because Module is a large struct. Passing by pointer avoids copying all the
// data every time we pass it to a function. The * means "pointer to".
//
// See: https://go.dev/tour/moretypes/1
// See: https://go.dev/doc/effective_go#pointers_vs_values
type Module struct {
	// --- Identity ---

	// Name is the directory name in the dotfiles repo (e.g., "nvim", "zsh").
	// This is used as the unique key to identify the module throughout the app.
	Name string

	// Icon is a Nerd Font glyph displayed next to the module name in the TUI.
	// Nerd Fonts patch developer-targeted fonts with many extra glyphs.
	// See: https://www.nerdfonts.com/cheat-sheet
	Icon string

	// Description is a brief one-line summary of what this module is.
	Description string

	// --- Categorization ---

	// Category groups related modules together in the UI sidebar.
	// Valid values: "Shell", "Editor", "Language", "DevOps", "Utility",
	// "Application", "Cloud", "Database", "AI".
	Category string

	// Website is the project's official website URL.
	Website string

	// Repo is the GitHub repository URL for the project.
	Repo string

	// --- Dependencies ---

	// Dependencies lists the Names of other modules in this repo that must
	// be installed before this one. For example, nvim depends on ["python",
	// "node", "go", "homebrew", "yazi"].
	//
	// This is a slice of strings where each string matches another Module's
	// Name field. The installer uses this to compute a topological install order.
	Dependencies []string

	// ExternalDeps lists tools needed from outside this dotfiles repo.
	// Unlike Dependencies (which reference other modules), these are system
	// packages that must be present for this module to work.
	ExternalDeps []ExternalDep

	// --- Installation ---

	// InstallCommands is an ordered list of shell commands to run during
	// installation. These are translated from the original install.sh scripts.
	// Each string is a single command (or pipeline) to execute via the shell.
	InstallCommands []string

	// UninstallCommands is an ordered list of shell commands to reverse
	// an installation. Not all modules have uninstall commands.
	UninstallCommands []string

	// StowEnabled indicates whether this module uses GNU Stow for symlinking
	// config files. Stow creates symlinks from ~/.config/<module> (or similar)
	// to the dotfiles repo directory.
	//
	// bool is Go's boolean type. Its zero value is false, so modules that
	// don't use stow don't need to set this field explicitly.
	// See: https://go.dev/ref/spec#Boolean_types
	StowEnabled bool

	// EstimatedTime is a rough human-readable install time (e.g., "30s", "2m").
	EstimatedTime string

	// EstimatedSize is a rough human-readable disk space estimate (e.g., "50MB").
	EstimatedSize string

	// CheckCommand is a shell command to verify the module is installed.
	// For example, "nvim --version" or "docker --version".
	// The exit code (0 = installed, non-zero = not installed) determines status.
	CheckCommand string

	// RequiresInput indicates whether the install process needs user interaction
	// (e.g., "press y to continue" prompts). This affects how the TUI handles
	// the installation — interactive installs may need special treatment.
	RequiresInput bool

	// ConfigOptions lists module-specific configuration choices that the user
	// can customize before installation.
	ConfigOptions []ConfigOption
}
