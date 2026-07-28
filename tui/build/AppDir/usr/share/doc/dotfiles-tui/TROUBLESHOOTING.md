# Troubleshooting

## App starts but the module list is empty

- Confirm `~/.config/dotfiles-tui/config.yaml` has a valid `modules_dir`.
- Or set `export DOTFILES_MODULES_DIR=/path/to/modules`.
- Each module needs a subdirectory with a `module.yaml` file.
- Invalid YAML is skipped; fix the file and restart.

## First-run path dialog keeps returning

A modules directory is required. Enter an absolute path or `~/something`. The
folder is created if it does not exist.

## Install fails with “module is required… Install it from the TUI first”

That module’s `install.sh` expects a dependency (e.g. homebrew, python) to
already be installed. Install the dependency module from the TUI first, or list
it under `dependencies:` in `module.yaml` so the orchestrator can order installs.

## Stow / symlink errors

- Ensure `stow` is installed (prereq screen).
- Config files must live under the module directory with the correct relative
  paths for your home layout (e.g. `.config/nvim/...`).
- Keep `install.sh`, `uninstall.sh`, and `module.yaml` in `.stow-local-ignore`.

## AppImage does not include modules

By design the AppImage is only the TUI + docs. Your modules live outside the
image on disk.

## After changing module.yaml nothing updates

Restart the TUI (or set the modules path again). The registry is loaded at
startup from disk.
