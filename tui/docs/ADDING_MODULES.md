# Adding New Modules

Each module is self-contained under `modules/<name>/`:

```
modules/starship/
  module.go      # metadata (registers with the TUI)
  install.sh     # install steps (optional)
  uninstall.sh   # uninstall steps (optional)
  .config/...    # files stowed to $HOME
  .stow-local-ignore
```

The TUI (`tui/` + `internal/`) is only the core app. Module add-ons live next to their scripts and configs.

## Steps

### 1. Create the package directory

```bash
mkdir -p modules/starship/.config/starship
# add config files...

cat > modules/starship/.stow-local-ignore << 'EOF'
install.sh
uninstall.sh
module.go
EOF

cat > modules/starship/install.sh << 'EOF'
#!/bin/bash
set -euo pipefail
curl -sS https://starship.rs/install.sh | sh -s -- --yes
EOF
chmod +x modules/starship/install.sh
```

### 2. Add `module.go` metadata

```go
package starship

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:          "starship", // must match directory name
		Icon:          "🚀",
		Description:   "Cross-shell prompt with starship.rs",
		Category:      "Shell", // used for sidebar filter (c)
		Website:       "https://starship.rs",
		Repo:          "https://github.com/starship/starship",
		Dependencies:  []string{},
		ExternalDeps:  nil,
		StowEnabled:   true,
		EstimatedTime: "~30s",
		EstimatedSize: "~15MB",
		CheckCommand:  "starship --version",
	})
}
```

`Category` groups modules in the UI. Press `c` in the TUI to filter by category.

Package name note: the Go package identifier cannot contain hyphens (`ai-tools` → `package aitools`). The directory name / `Name` field stay hyphenated.

### 3. Register the package with the TUI

Add a blank import in [`tui/modules_register.go`](../modules_register.go):

```go
_ "github.com/issafalcon/dotfiles-tui/modules/starship"
```

(Or regenerate that file by listing every `modules/*/module.go`.)

### 4. Verify

```bash
cd tui
make run
```

Press `c` to open the category picker, `i` to install, `r` on confirm to review `install.sh`.

## Field reference

| Field | Required | Notes |
|-------|----------|-------|
| `Name` | yes | Directory under `modules/` |
| `Category` | yes | Shell, Editor, Language, DevOps, Cloud, Database, Utility, Application, AI |
| `Description` | yes | Sidebar one-liner |
| `CheckCommand` | yes | Exit 0 = installed |
| `Dependencies` | | Other module names |
| `ExternalDeps` | | System tools not managed as modules |
| `StowEnabled` | | Symlink configs via GNU Stow |

Install/uninstall logic stays in the shell scripts — do not put command lists in Go.
