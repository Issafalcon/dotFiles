# Adding New Modules

This guide explains how to add a new dotfile module to the DotFiles TUI app.

## Overview

Each module in the app corresponds to a directory in the dotfiles repository root. The TUI manages:
1. **Installing** the module's dependencies (via `install.sh` commands)
2. **Stowing** the module's configuration files to `$HOME` (via GNU Stow)
3. **Tracking** which modules are installed (via `~/.dotFileModules`)

## Step-by-Step

### 1. Create the Dotfile Module Directory

First, create the module directory in the dotfiles repo root with its configuration files:

```bash
# Example: adding a "starship" prompt module
mkdir starship

# Add config files that should be symlinked to $HOME
mkdir -p starship/.config/starship
# Place starship.toml in starship/.config/starship/starship.toml

# Create .stow-local-ignore to exclude non-config files
cat > starship/.stow-local-ignore << 'EOF'
install.sh
EOF

# Create install.sh for installing the tool itself
cat > starship/install.sh << 'EOF'
#!/bin/bash
curl -sS https://starship.rs/install.sh | sh
EOF
chmod +x starship/install.sh
```

### 2. Create the Module Definition File

Create a new Go file in `tui/internal/module/modules/`. You can either:
- Add to an existing category file (e.g., `shell.go` for shell-related tools)
- Create a new file for the module

```go
// File: tui/internal/module/modules/starship.go
package modules

import (
    "github.com/issafalcon/dotfiles-tui/internal/module"
)

// init() is a special Go function that runs automatically when the package
// is imported. We use it to register modules with the global registry.
// Each module file's init() runs once at program startup.
//
// See: https://go.dev/doc/effective_go#init
func init() {
    module.DefaultRegistry.Register(&module.Module{
        // Name MUST match the directory name in the dotfiles repo root.
        // This is used by GNU Stow to find the module's config files.
        Name: "starship",

        // Icon: A Nerd Font glyph. Find icons at https://www.nerdfonts.com/cheat-sheet
        Icon: "🚀",

        // Description: One-line summary shown in the sidebar module list.
        Description: "Cross-shell prompt with starship.rs",

        // Category: Used for grouping in the UI.
        // Valid categories: Shell, Editor, Language, DevOps, Cloud,
        //                   Database, Utility, Application, AI
        Category: "Shell",

        // Website: The project's homepage (opened with 'o' key).
        Website: "https://starship.rs",

        // Repo: GitHub repository URL.
        Repo: "https://github.com/starship/starship",

        // Dependencies: Names of OTHER modules that must be installed first.
        // These must match the Name field of other registered modules.
        // The install orchestrator will install these before this module.
        Dependencies: []string{},

        // ExternalDeps: System tools needed that aren't managed as modules.
        // Each has a check command, install command, and method.
        ExternalDeps: []module.ExternalDep{
            {
                Name:           "curl",
                CheckCommand:   "curl --version",
                InstallCommand: "sudo apt install -y curl",
                InstallMethod:  "apt",
            },
        },

        // InstallCommands: Shell commands run to install the tool.
        // These should mirror what's in the module's install.sh file.
        // Commands are run sequentially in order.
        InstallCommands: []string{
            "curl -sS https://starship.rs/install.sh | sh -s -- --yes",
        },

        // UninstallCommands: Shell commands to remove the tool.
        // Leave empty if uninstall isn't supported.
        UninstallCommands: []string{
            "sh -c 'rm -f $(which starship)'",
        },

        // StowEnabled: Set to true if this module has config files
        // that should be symlinked to $HOME via GNU Stow.
        StowEnabled: true,

        // EstimatedTime: Rough estimate shown in the UI.
        EstimatedTime: "~30s",

        // EstimatedSize: Rough disk space estimate.
        EstimatedSize: "~15MB",

        // CheckCommand: Command used to verify installation succeeded.
        // The app runs this after install and checks the exit code.
        CheckCommand: "starship --version",

        // RequiresInput: Set to true if the install commands need
        // user interaction (e.g., confirmation prompts).
        // The UI will handle input forwarding if this is true.
        RequiresInput: false,

        // ConfigOptions: Module-specific settings the user can configure
        // in the Configuration tab. Leave empty for simple modules.
        ConfigOptions: []module.ConfigOption{},
    })
}
```

### 3. Import the Module Package

The modules are auto-registered via `init()` functions. The module package imports all module definition files using a blank import. If you created a new file in the existing `modules/` directory, it's already included.

If you created a new subdirectory, add a blank import in the appropriate place:

```go
import _ "github.com/issafalcon/dotfiles-tui/internal/module/modules"
```

### 4. Verify

```bash
cd tui
go build ./...   # Compile check
make run          # Visual check — your module should appear in the sidebar
```

## Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | `string` | ✅ | Directory name in dotfiles repo |
| `Icon` | `string` | ✅ | Nerd Font glyph or emoji |
| `Description` | `string` | ✅ | One-line summary |
| `Category` | `string` | ✅ | Grouping category |
| `Website` | `string` | | Project homepage URL |
| `Repo` | `string` | | GitHub repo URL |
| `Dependencies` | `[]string` | | Other module names |
| `ExternalDeps` | `[]ExternalDep` | | System tool dependencies |
| `InstallCommands` | `[]string` | ✅ | Shell commands to install |
| `UninstallCommands` | `[]string` | | Shell commands to uninstall |
| `StowEnabled` | `bool` | | Whether to use GNU Stow |
| `EstimatedTime` | `string` | | e.g., "~2min" |
| `EstimatedSize` | `string` | | e.g., "~50MB" |
| `CheckCommand` | `string` | ✅ | Verify install command |
| `RequiresInput` | `bool` | | Needs user interaction? |
| `ConfigOptions` | `[]ConfigOption` | | User-configurable options |

## Tips

- **Mirror install.sh**: Your `InstallCommands` should match the existing `install.sh` as closely as possible
- **Dependencies matter**: If your module needs `homebrew`, list it in `Dependencies`
- **Check commands**: Use `<tool> --version` or `command -v <tool>` for verification
- **Icons**: Browse [Nerd Fonts Cheat Sheet](https://www.nerdfonts.com/cheat-sheet) for appropriate icons
- **Testing**: After adding a module, test the full install flow with `make run`
