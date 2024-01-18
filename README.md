# Issafalcon dotfiles

> Modular installation of terminal, and terminal tools with my personal config files
> Inspired by [`caarlos0 dotFiles setup`](https://github.com/caarlos0/dotfiles)
>
> DISCLOSURE: Most of the contents have only been tested using Ubuntu 22.04 on native Linux machine and WSL
>             They are also an ongoing WIP and highly personalised, so please review the module code before you install to make sure it fits your needs

The goals of my dotFiles are as follows:
 1. Replace default shell with zsh, adding useful plugins without compromising speed
 2. Create modular installation options for remaining dotfiles
 3. Allow customization and extensibility

## Installation

### Setup and prerequisites intallation

The following will install some prerequisite files onto your machine (requires `sudo`)
e.g. git, curl, wget etc.

It will quickly replace the default shell with zsh and add plugins using zinit.

The default theme is powerline10k (you can change this)

These prerequisites may change as I evolve this repo.

> IMPORTANT: Clone to the ~/dotFiles directory as this will become the default stow directory (this can of course be modified)

```console
$ git clone https://github.com/Issafalcon/dotFiles.git ~/dotFiles
$ cd ~/dotFiles
$ ./prerequisites.sh 
$ zsh
```

Zinit plugins will be installed and you will need to restart your terminal to be taken to
powerline10k configuration wizard.

It is recommended that you install a NerdFont compatible font prior to setting up powerline10k.

### Module Installation

Following the `prerequisites.sh` script run, modules can be installed via two scripts in the root directory:
- `./bootstrap.sh`
- `./bootstrap_bulk.sh`

Module names match the names of the top level subdirectories in the repo (except `docs` and `.git`)

The recommended order of modules would be:
  1. `fzf` - The `z` function (fzf search on previously visited directories) relies on this being installed
  2. `homebrew` - Required for installation of some of the other modules
  3. `libsecret` - Used for storing git credentials in WSL
  4. `git` - Adds `delta` which provides prettier output for git commands
  5. `tmux` - Terminal multiplexer
  6. `ranger` - Terminal file explorer
  7. `lazygit` - TUI for git

#### Individual Modules

To "install" a module run the following:

```console
$ ./bootstrap.sh -m <MODULE_NAME> -i
```
> Or, if you don't want to install the actual dependencies for the module (you may have them installed already)

```console
$ ./bootstrap.sh -m <MODULE_NAME>
```
#### Bulk Install Modules

To install multiple modules at once (or all of them), run the following:

```console
$ ./bootstrap_bulk.sh -i [...MODULE_NAMES] # Where <MODULE_NAMES> is a space separated list of modules
```
> Or, if you don't want to install the actual dependencies for the module (you may have them installed already)

```console
$ ./bootstrap.sh -m [...MODULE_NAMES]
```

### Neovim Setup

### MySql Setup

This module folder contains a `docker-compose.yaml` file and associated `.env` and SQL script files to startup a containerized MySql server.

## Further help:

- [Opinionated Terminal Setup for WSL2 on Windows](/docs/WSL2.md)
- [Personalize your configs](/docs/PERSONALIZATION.md)
- [Understand how it works](/docs/DESIGN.md)

## Contributing

At the moment, I am not accepting PRs, but please feel free to open issues or add suggestions for improvements,
so that setup can be made more accessible.

## Feature Roadmap 🌌:
- [x] Replace all occurences of ~/repos with $PROJECTS variable
- [x] Replace Homebrew installations with apt package manager (and remove Homebrew from deps) (NOTE: Lazygit and some language servers / formatters still require homebrew)
- [ ] Add uninstall scripts for modules
  - [ ] zsh
  - [ ] node
  - [ ] python
  - [ ] etc...
- [ ] Add ability to load custom .zsh settings
- [ ] Add detailed design guide to how it works
- [x] Tidy up prerequisites to make as minimal as possible
- [ ] Add silent install for modules
  - [ ] zsh
  - [ ] node
  - [ ] python
  - [ ] etc...

