# Installation

## Quick start

```console
git clone https://github.com/Issafalcon/dotFiles.git ~/dotFiles
cd ~/dotFiles/tui
go run .
# or: make run / make build && ./build/dotfiles-tui
```

The TUI checks for **git**, **stow**, and **curl** on startup and can install any that are missing. Install modules from the dashboard (`i` to install, `d` to uninstall). After shell-related installs, open a new terminal so PATH / `.zshrc` changes apply.

Optional: set `DOTFILES_DIR` if the repo is not auto-detected from the working directory.

## Prerequisites

| Tool | Why | Install |
|------|-----|---------|
| Go 1.22+ | Build/run the TUI | https://go.dev/dl/ |
| `git` | Clone / version control | `sudo apt install git` |
| `stow` | Symlink module configs | `sudo apt install stow` |
| `curl` | Many module install scripts | `sudo apt install curl` |

Everything else (zsh, node, build tools, …) is installed by selecting the corresponding module in the TUI.

## Recommended first modules

1. `zsh` — shell + zinit
2. `fzf` — fuzzy finder / `z`
3. `homebrew` — needed by some modules
4. `git`, `tmux`, `nvim` — as needed

## Building

```console
cd tui
make build          # → build/dotfiles-tui
make install        # → /usr/local/bin/dotfiles-tui
```
