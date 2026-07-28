// This file provides post-installation verification for the DotFiles TUI.
//
// After installing a module, we verify that the expected binary/command is
// available and report its version. This gives the user confidence that the
// installation succeeded and the tool is properly set up.
//
// Key Go concepts used here:
//   - tea.Cmd pattern: Wrapping verification in a Bubble Tea command
//   - Struct literals: Creating structs with named fields
//   - Error handling: Go's explicit error checking pattern
package installer

import (
	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/dotfiles-tui/internal/utils"
)

// VerifyMsg is a Bubble Tea message sent after post-install verification completes.
//
// The Update() handler receives this and can update the UI to show a checkmark
// (if installed) or a warning (if verification failed).
//
// In Go, struct types are defined with the type keyword. Structs are value types,
// meaning they are copied when assigned or passed to functions (unless you use pointers).
// See: https://go.dev/ref/spec#Struct_types
// See: https://go.dev/tour/moretypes/2
type VerifyMsg struct {
	// ModuleName identifies which module was verified.
	ModuleName string

	// Installed is true if the check command was found in PATH.
	Installed bool

	// Version is the version string reported by the command (empty if not installed).
	Version string
}

// VerifyInstall returns a tea.Cmd that checks if a binary is available after install.
//
// This follows the Bubble Tea command pattern:
//
//  1. VerifyInstall is called with the command to check (e.g., "nvim").
//  2. It returns a tea.Cmd (a function that returns tea.Msg).
//  3. Bubble Tea executes this function in a goroutine.
//  4. The function checks if the command exists and gets its version.
//  5. It returns a VerifyMsg that flows into Update().
//
// Parameters:
//   - moduleName: The name of the module being verified (for the message).
//   - checkCommand: The command to check for (e.g., "git", "nvim", "stow").
//
// Returns:
//   - tea.Cmd: A command function that Bubble Tea will execute asynchronously.
//
// Usage in a Bubble Tea Update() handler:
//
//	case InstallCompleteMsg:
//	    if msg.Success {
//	        return m, installer.VerifyInstall("nvim", "nvim")
//	    }
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
func VerifyInstall(moduleName string, checkCommand string) tea.Cmd {
	// Return an anonymous function (closure) that captures moduleName and checkCommand.
	//
	// This is the standard Bubble Tea pattern for creating commands. The closure
	// runs in a separate goroutine managed by the Bubble Tea runtime. It must not
	// modify the model directly — it communicates only by returning a message.
	//
	// Anonymous functions in Go are also called "function literals".
	// See: https://go.dev/ref/spec#Function_literals
	return func() tea.Msg {
		// Check if the command is available in PATH.
		installed := utils.IsCommandAvailable(checkCommand)

		// Initialize version as empty string.
		// In Go, var declares a variable with its zero value.
		// The zero value for string is "" (empty string).
		// See: https://go.dev/ref/spec#The_zero_value
		var version string

		if installed {
			// Try to get the version string. If it fails, we still report
			// the command as installed — we just won't have version info.
			//
			// The if-init statement is a Go idiom that limits the scope of
			// variables. The 'v' and 'err' variables are only accessible
			// within this if/else block.
			// See: https://go.dev/doc/effective_go#if
			if v, err := utils.GetCommandVersion(checkCommand); err == nil {
				version = v
			}
		}

		// Return a VerifyMsg using a struct literal with named fields.
		// Named fields make the code self-documenting and order-independent.
		// See: https://go.dev/ref/spec#Composite_literals
		return VerifyMsg{
			ModuleName: moduleName,
			Installed:  installed,
			Version:    version,
		}
	}
}
