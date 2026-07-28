// Package module defines the core data structures for dotfile modules.
//
// A "module" in this context is a single tool or application whose configuration
// is managed by this dotfiles repository. Each module has metadata (name, icon,
// description) and dependency information for ordering installations.
// Install/uninstall steps live in modules/<name>/install.sh and uninstall.sh;
// the TUI executes those scripts rather than embedding shell commands in Go.
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
	Name           string `yaml:"name"`            // Human-readable name of the external tool.
	CheckCommand   string `yaml:"check_command"`   // Shell command to verify the tool is installed.
	InstallCommand string `yaml:"install_command"` // Shell command to install the tool if missing.
	InstallMethod  string `yaml:"install_method"`  // Package manager: apt, brew, cargo, npm, pip, or curl.
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
	Name        string   `yaml:"name"`        // Short identifier for this option (e.g., "dotnet_version").
	Description string   `yaml:"description"` // Human-readable description shown to the user.
	Default     string   `yaml:"default"`     // Default value if the user doesn't choose.
	Choices     []string `yaml:"choices"`     // Valid values the user can pick from. Empty means freeform.
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

	// Name is the directory name under the modules path (e.g., "nvim", "zsh").
	Name string `yaml:"name"`

	// Icon is a Nerd Font glyph displayed next to the module name in the TUI.
	Icon string `yaml:"icon"`

	// Description is a brief one-line summary of what this module is.
	Description string `yaml:"description"`

	// --- Categorization ---

	// Category groups related modules together in the UI sidebar.
	Category string `yaml:"category"`

	// Website is the project's official website URL.
	Website string `yaml:"website"`

	// Repo is the GitHub repository URL for the project.
	Repo string `yaml:"repo"`

	// --- Dependencies ---

	Dependencies []string      `yaml:"dependencies"`
	ExternalDeps []ExternalDep `yaml:"external_deps"`

	// --- Installation ---

	StowEnabled   bool           `yaml:"stow_enabled"`
	EstimatedTime string         `yaml:"estimated_time"`
	EstimatedSize string         `yaml:"estimated_size"`
	CheckCommand  string         `yaml:"check_command"`
	RequiresInput bool           `yaml:"requires_input"`
	ConfigOptions []ConfigOption `yaml:"config_options"`
}
