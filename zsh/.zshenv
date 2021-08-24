# Dotfile reference
export DOTFILES=$HOME/dotFiles

# your project folder that we can `c [tab]` to
export PROJECTS="$HOME/repos"

export EDITOR="nvim"

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


