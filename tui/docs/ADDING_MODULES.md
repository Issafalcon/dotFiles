# Adding Modules

See the embedded guide (same content as shipped in the AppImage and shown with `H` in the TUI):

The canonical copy lives at [`internal/docs/ADDING_MODULES.md`](../../internal/docs/ADDING_MODULES.md).

Quick summary:

1. Create `$MODULES_DIR/<name>/module.yaml` (+ optional `install.sh` / `uninstall.sh` + stow tree).
2. Restart the TUI (or set `DOTFILES_MODULES_DIR`) — no recompile.
3. Press `c` to filter by `category`, `i` to install, `r` to review the script.
