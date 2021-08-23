# Dotfile reference
export DOTFILES=$HOME/.dotfiles

# your project folder that we can `c [tab]` to
export PROJECTS="$HOME/repos"

export EDITOR="nvim"
export VISUAL="nvim"
export LSCOLORS='exfxcxdxbxegedabagacad'
export CLICOLOR=true

export HISTFILE=~/.zsh_history
export HISTSIZE=10000
export SAVEHIST=10000

export MANPATH="/usr/local/man:$MANPATH"

# nvm
export NVM_DIR="$HOME/.nvm"


# TODO: Move the scripts into a custom scripts folder inside dotFiles
if [ -d "$HOME/repos/scripts" ]; then
  export PATH=$HOME/repos/scripts/bash/utilities:$PATH
  export PATH=$HOME/repos/scripts/bash/terminal:$PATH
fi

if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # set DISPLAY variable to the IP automatically assigned to WSL2
  export DISPLAY=$(route.exe print | grep 0.0.0.0 | head -1 | awk '{print $4}'):0.0

  # Used for vagrant - Enables vagrant use from within WSL2
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
  export PATH="$PATH:/mnt/c/Program Files/Oracle/VirtualBox"
fi

