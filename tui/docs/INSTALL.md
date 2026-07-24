# Installation

## AppImage (recommended for end users)

1. Download the latest `dotfiles-tui-*.AppImage` release asset.
2. `chmod +x dotfiles-tui-*.AppImage && ./dotfiles-tui-*.AppImage`
3. Install any missing **git / stow / curl** from the prereq screen.
4. Enter the path to your **modules** directory (created if missing). The AppImage ships with **no** modules — you bring your own.

Optional:

```bash
export DOTFILES_MODULES_DIR=~/my-dotfiles/modules
./dotfiles-tui-*.AppImage
```

Config is saved to `~/.config/dotfiles-tui/config.yaml`.

## From source (developers)

```bash
git clone <this-repo> ~/dotFiles
cd ~/dotFiles/tui
make run    # auto-detects ../modules
# or: make build && ./build/dotfiles-tui
```

## Building an AppImage

Requires `appimagetool` on `PATH`:

```bash
cd tui
make appimage
# → build/dotfiles-tui-<version>-x86_64.AppImage
```

The image contains the TUI binary, desktop entry, icon, and docs under
`usr/share/doc/dotfiles-tui/` — not a modules tree.

## Keys

| Key | Action |
|-----|--------|
| `c` | Filter by category |
| `H` | Adding-modules guide |
| `i` | Install |
| `r` | Review install.sh on confirm |
| `d` | Uninstall |
| `?` | Help |

See [ADDING_MODULES.md](./ADDING_MODULES.md) (also embedded; press `H`).
