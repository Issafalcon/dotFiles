# Dotfile reference
export DOTFILES=$HOME/dotFiles

# your project folder that we can `c [tab]` to
export PROJECTS="$HOME"/repos

# Custom terminfo allowing for italic display
export TERM=xterm-256color-italic

export EDITOR="nvim"

export LSCOLORS='exfxcxdxbxegedabagacad'
export CLICOLOR=true

export HISTFILE=~/.zsh_history
export HISTSIZE=10000
export SAVEHIST=10000

export MANPATH="/usr/local/man:$MANPATH"

# nvm
export NVM_DIR="$HOME/.nvm"

# Allow WSL gui apps to load up in Windows X Server
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # set DISPLAY variable to the IP automatically assigned to WSL2
  export DISPLAY="`grep nameserver /etc/resolv.conf | sed 's/nameserver //'`:0"
  path+=("/mnt/c/Program Files/Oracle/VirtualBox")
  # Used for vagrant - Enables vagrant use from within WSL2
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
fi
