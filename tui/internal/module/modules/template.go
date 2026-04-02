// Package modules contains all individual module definitions for the DotFiles TUI.
//
// This is a TEMPLATE FILE — a heavily documented example showing how to create
// a new module definition. Copy this pattern when adding a new module.
//
// # How Module Registration Works
//
// Each module file uses Go's init() function to register its module(s) with
// the global DefaultRegistry. Here's the lifecycle:
//
//  1. Go compiles all files in the modules package
//  2. At program startup, Go runs init() functions for every imported package
//  3. Each init() function calls module.DefaultRegistry.Register(...)
//  4. By the time main() runs, all modules are registered and ready to use
//
// # The init() Function
//
// init() is a special function in Go. You don't call it — Go calls it
// automatically when the package is loaded. A single file can have multiple
// init() functions, and a package can have init() in every file. They all
// run in the order the source files are compiled (alphabetically by filename).
//
// See: https://go.dev/doc/effective_go#init
// See: https://go.dev/ref/spec#Package_initialization
//
// # Import Side Effects
//
// For module registration to work, the modules package must be imported
// somewhere in the application. Even though no functions from this package
// are called directly, the import triggers all init() functions to run.
// Go uses the blank identifier import for this:
//
//	import _ "github.com/issafalcon/dotfiles-tui/internal/module/modules"
//
// The underscore (_) tells Go: "import this package for its side effects only"
// (i.e., run its init functions). Without this import, none of the modules
// would be registered!
//
// See: https://go.dev/doc/effective_go#blank_import
//
// # Naming Conventions
//
// - File names: use the module name or a group name (e.g., "cloud.go" for
//   aws/azure/google-cloud, "languages.go" for go/rust/python/node/lua/cpp/dotnet)
// - Function names: init() for registration (Go convention)
// - Package name: "modules" (plural, since it contains many modules)
//
// See: https://go.dev/doc/effective_go#names
// See: https://go.dev/blog/package-names
package modules

// ----- TEMPLATE: Example Module Definition -----
//
// Below is a fully commented example of a module definition.
// Uncomment and modify it to create a new module.
//
// import "github.com/issafalcon/dotfiles-tui/internal/module"
//
// func init() {
//     // Register this module with the global DefaultRegistry.
//     // The & operator creates a pointer to the Module struct literal.
//     // See: https://go.dev/tour/moretypes/1
//     module.DefaultRegistry.Register(&module.Module{
//
//         // --- Identity Fields ---
//
//         // Name: The directory name in the dotfiles repo.
//         // This MUST match the actual directory name (e.g., "nvim", "zsh").
//         // It's used as the unique key in the registry.
//         Name: "example-tool",
//
//         // Icon: A Nerd Font glyph for the TUI display.
//         // Find glyphs at: https://www.nerdfonts.com/cheat-sheet
//         // Use the actual Unicode character, not an escape sequence.
//         Icon: "",
//
//         // Description: A brief one-line summary shown in the module list.
//         Description: "An example tool for demonstration purposes",
//
//         // --- Categorization ---
//
//         // Category: The grouping for the sidebar. Must be one of:
//         // "Shell", "Editor", "Language", "DevOps", "Utility",
//         // "Application", "Cloud", "Database", "AI"
//         Category: "Utility",
//
//         // Website: The project's official website.
//         Website: "https://example.com",
//
//         // Repo: The GitHub repository URL.
//         Repo: "https://github.com/example/tool",
//
//         // --- Dependencies ---
//
//         // Dependencies: Other module Names that must be installed first.
//         // Use nil or omit if there are no dependencies.
//         // Each string must match an existing module's Name field.
//         Dependencies: []string{"python", "node"},
//
//         // ExternalDeps: System tools needed that aren't other modules.
//         // Each ExternalDep specifies how to check for and install the tool.
//         ExternalDeps: []module.ExternalDep{
//             {
//                 Name:           "curl",
//                 CheckCommand:   "curl --version",
//                 InstallCommand: "sudo apt-get install -y curl",
//                 InstallMethod:  "apt",
//             },
//         },
//
//         // --- Installation ---
//
//         // InstallCommands: Shell commands to run, in order.
//         // These come from the original install.sh scripts.
//         // Each string is passed to the shell for execution.
//         InstallCommands: []string{
//             "sudo apt-get update",
//             "sudo apt-get install -y example-tool",
//         },
//
//         // UninstallCommands: Shell commands to reverse the installation.
//         // Can be nil if uninstall isn't supported.
//         UninstallCommands: []string{
//             "sudo apt-get remove -y example-tool",
//         },
//
//         // StowEnabled: true if this module's config directory should be
//         // symlinked by GNU Stow. Most modules with config files use stow.
//         StowEnabled: true,
//
//         // EstimatedTime: Rough install time for user expectations.
//         EstimatedTime: "30s",
//
//         // EstimatedSize: Rough disk space usage.
//         EstimatedSize: "10MB",
//
//         // CheckCommand: Shell command to verify installation.
//         // Exit code 0 = installed, non-zero = not installed.
//         CheckCommand: "example-tool --version",
//
//         // RequiresInput: Set to true if the install script has prompts
//         // (e.g., "Press y to continue"). The TUI handles these specially.
//         RequiresInput: false,
//
//         // ConfigOptions: User-configurable choices for this module.
//         // Can be nil if there are no configuration options.
//         ConfigOptions: []module.ConfigOption{
//             {
//                 Name:        "theme",
//                 Description: "Color theme to use",
//                 Default:     "dark",
//                 Choices:     []string{"dark", "light", "auto"},
//             },
//         },
//     })
// }
