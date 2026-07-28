# Adding Modules

The DotFiles TUI discovers modules at **runtime** from a directory you choose
(first launch, or `DOTFILES_MODULES_DIR` / `~/.config/dotfiles-tui/config.yaml`).

No rebuild or AppImage update is required to add modules.

## Layout

```
$MODULES_DIR/
  my-tool/
    module.yaml           # required metadata
    install.sh            # optional
    uninstall.sh          # optional
    .config/...           # files for GNU Stow
    .stow-local-ignore    # ignore scripts + module.yaml
```

## module.yaml

```yaml
name: my-tool                 # should match the directory name
icon: ""
description: One-line summary
category: Utility             # Shell, Editor, Language, DevOps, Cloud, Database, Utility, Application, AI
website: https://example.com
repo: https://github.com/example/tool
dependencies: []              # other module directory names
external_deps:
  - name: curl
    check_command: curl --version
    install_command: sudo apt-get install -y curl
    install_method: apt
stow_enabled: true
estimated_time: 30s
estimated_size: 10MB
check_command: my-tool --version
requires_input: false
```

The directory name is the source of truth for `name` if they disagree.

## Scripts

- `install.sh` — run by the TUI before stow (when present)
- `uninstall.sh` — run before unstow (when present)

Ignore both (and `module.yaml`) in `.stow-local-ignore` so they are not linked into `$HOME`.

## AppImage users

1. Download and run the AppImage (`chmod +x` first).
2. Pass the prereq screen (git, stow, curl).
3. Enter your modules directory path (created if missing).
4. Add modules under that path; restart or re-open the app to pick up new `module.yaml` files.
5. Press `H` in the dashboard to re-read this guide; `c` to filter by category; `i` to install; `r` on confirm to review `install.sh`.

## Developing from this repository

If you run the TUI from the repo checkout, a sibling `modules/` directory is
auto-detected. You can still override with `DOTFILES_MODULES_DIR`.
