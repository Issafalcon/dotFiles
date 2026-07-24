# Issafalcon dotfiles

> Modular installation of terminal tools and configs, managed by a Go TUI
> Inspired by [`caarlos0 dotFiles setup`](https://github.com/caarlos0/dotfiles)
>
> DISCLOSURE: Most of the contents have only been tested using Ubuntu 22.04 on native Linux and WSL.
> They are an ongoing WIP and highly personalised — review module scripts before installing.

## Goals

1. Replace the default shell with zsh and useful plugins without compromising speed
2. Modular install/uninstall of each tool’s configs via GNU Stow
3. Drive everything from an interactive TUI (no separate bootstrap scripts)

## Layout

```
go.mod / internal/   # Go module root + TUI core packages
tui/                 # main package, Makefile, docs
modules/<name>/      # install.sh, uninstall.sh, module.go, stow configs
```

Press `c` in the TUI to filter the sidebar by category.
## Installation

Clone the repo, then run the TUI:

```console
git clone https://github.com/Issafalcon/dotFiles.git ~/dotFiles
cd ~/dotFiles/tui
go run .
```

On first launch the app checks for `git`, `stow`, and `curl`. Install any missing tools from that screen, then use the dashboard to install modules (`i`), uninstall (`d`), or review an install script (`r` on the confirm dialog) before confirming.

After installing `zsh` (or other PATH-changing modules), open a new terminal so shell config takes effect.

It is recommended that you install a Nerd Font before setting up Powerlevel10k via the zsh module.

### Suggested install order

1. `zsh`
2. `fzf`
3. `homebrew` (required by some modules)
4. `libsecret` / `git` / `tmux` / `nvim` as needed

## Further help

- [Opinionated Terminal Setup for WSL2 on Windows](/docs/WSL2.md)
- [Personalize your configs](/docs/PERSONALIZATION.md)
- [Understand how it works](/docs/DESIGN.md)
- [TUI install notes](/tui/docs/INSTALL.md)
- [Adding modules](/tui/docs/ADDING_MODULES.md)

## Contributing

At the moment I am not accepting PRs, but feel free to open issues or suggestions.
