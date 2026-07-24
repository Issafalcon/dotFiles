# Issafalcon dotfiles

> Modular terminal tools and configs, managed by a Go TUI (AppImage-friendly).
>
> DISCLOSURE: Tested mainly on Ubuntu 22.04 (native / WSL). Review module scripts before installing.

## Quick start (AppImage)

1. Download the AppImage, `chmod +x`, run it.
2. Satisfy prereqs (git, stow, curl) if prompted.
3. Point the app at a **modules** directory you own (empty is fine).
4. Add modules as folders with `module.yaml` (+ scripts/configs). See docs in-app (`H`) or [internal/docs/ADDING_MODULES.md](internal/docs/ADDING_MODULES.md).

The AppImage does **not** bundle modules — only the TUI and documentation.

## Layout

```
go.mod / internal/     # TUI core
tui/                   # main, Makefile, AppImage assets
modules/<name>/        # your packages: module.yaml, install.sh, stow files
```

## From source

```console
git clone <repo> ~/dotFiles
cd ~/dotFiles/tui && make run
```

## Further help

- [Install / AppImage](tui/docs/INSTALL.md)
- [Adding modules](internal/docs/ADDING_MODULES.md)
- [Troubleshooting](internal/docs/TROUBLESHOOTING.md)
