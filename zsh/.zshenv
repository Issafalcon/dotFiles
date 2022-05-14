# Dotfile reference
export DOTFILES=$HOME/dotFiles

# your project folder that we can `c [tab]` to
export PROJECTS="$HOME"/repos

export TERM=xterm-256color

export EDITOR="nvim"

export LSCOLORS='exfxcxdxbxegedabagacad'
export CLICOLOR=true

export HISTFILE=~/.zsh_history
export HISTSIZE=10000
export SAVEHIST=10000

export MANPATH="/usr/local/man:$MANPATH"

# nvm
export NVM_DIR="$HOME/.nvm"

# Setting local (TODO: Uncomment the en_GB.UTF-8 in /etc/locale.gen after running 'sudo locale'. Then run sudo locale-gen after)
export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# Allow WSL gui apps to load up in Windows X Server
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # set DISPLAY variable to the IP automatically assigned to WSL2
  export DISPLAY=export DISPLAY=$(ip route  | awk '/default via / {print $3; exit}' 2>/dev/null):0
  path+=("/mnt/c/Program Files/Oracle/VirtualBox")
  # Used for vagrant - Enables vagrant use from within WSL2
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
fi

[ -s "$HOME/.zshenv_local" ] && source "$HOME/.zshenv_local"
