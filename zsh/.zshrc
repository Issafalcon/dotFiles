#!/bin/zsh
# uncomment this and the last line for zprof info
# zmodload zsh/zprof

export DOTFILES=$HOME/.dotfiles

# your project folder that we can `c [tab]` to
export PROJECTS="$HOME/repos"

export EDITOR="nvim"
export VISUAL="nvim"

# Check what modules have been 'installed'
MODULES = $(cat $HOME/.dotFileModules)

export LSCOLORS='exfxcxdxbxegedabagacad'
export CLICOLOR=true

fpath=($DOTFILES/functions $fpath)

autoload -U "$DOTFILES"/functions/*(:t)
autoload -U up-line-or-beginning-search
autoload -U down-line-or-beginning-search
autoload -U edit-command-line

HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000

#Configurations options
  # don't nice background tasks
 setopt NO_BG_NICE
 setopt NO_HUP
 setopt NO_BEEP
 # allow functions to have local options
 setopt LOCAL_OPTIONS
 # allow functions to have local traps
 setopt LOCAL_TRAPS
 # share history between sessions ???
 setopt SHARE_HISTORY
 # add timestamps to history
 setopt EXTENDED_HISTORY
 setopt PROMPT_SUBST
 setopt CORRECT
 setopt COMPLETE_IN_WORD
 # adds history
 setopt APPEND_HISTORY
 # adds history incrementally and share it across sessions
 setopt INC_APPEND_HISTORY
 setopt SHARE_HISTORY
 # don't record dupes in history
 setopt HIST_IGNORE_ALL_DUPS
 setopt HIST_REDUCE_BLANKS
 setopt HIST_IGNORE_DUPS
 setopt HIST_IGNORE_SPACE
 setopt HIST_VERIFY
 setopt HIST_EXPIRE_DUPS_FIRST
 # dont ask for confirmation in rm globs*
 setopt RM_STAR_SILENT

 zle -N up-line-or-beginning-search
 zle -N down-line-or-beginning-search
 zle -N edit-command-line

 # fuzzy find: start to type
 bindkey "$terminfo[kcuu1]" up-line-or-beginning-search
 bindkey "$terminfo[kcud1]" down-line-or-beginning-search
 bindkey "$terminfo[cuu1]" up-line-or-beginning-search
 bindkey "$terminfo[cud1]" down-line-or-beginning-search

 # backward and forward word with option+left/right
 bindkey '^[^[[D' backward-word
 bindkey '^[b' backward-word
 bindkey '^[^[[C' forward-word
 bindkey '^[f' forward-word

 # to to the beggining/end of line with fn+left/right or home/end
 bindkey "${terminfo[khome]}" beginning-of-line
 bindkey '^[[H' beginning-of-line
 bindkey "${terminfo[kend]}" end-of-line
 bindkey '^[[F' end-of-line

 # delete char with backspaces and delete
 bindkey '^[[3~' delete-char
 bindkey '^?' backward-delete-char

 # delete word with ctrl+backspace
 bindkey '^[[3;5~' backward-delete-word
 # bindkey '^[[3~' backward-delete-word

 # edit command line in $EDITOR
 bindkey '^e' edit-command-line

 # search history with fzf if installed, default otherwise
 if test -d /usr/local/opt/fzf/shell; then
   # shellcheck disable=SC1091
     . /usr/local/opt/fzf/shell/key-bindings.zsh
     else
       bindkey '^R' history-incremental-search-backward
       fi

# Load module path files
for module in "${MODULES}"; do
  source "$DOTFILES/$module/path.zsh"
done


autoload -Uz promptinit
promptinit
prompt adam1

setopt histignorealldups sharehistory

# Use emacs keybindings even if our EDITOR is set to vi
bindkey -e

# Keep 1000 lines of history within the shell and save it to ~/.zsh_history:
HISTSIZE=1000
SAVEHIST=1000
HISTFILE=~/.zsh_history

# Use modern completion system
autoload -Uz compinit
compinit

zstyle ':completion:*' auto-description 'specify: %d'
zstyle ':completion:*' completer _expand _complete _correct _approximate
zstyle ':completion:*' format 'Completing %d'
zstyle ':completion:*' group-name ''
zstyle ':completion:*' menu select=2
eval "$(dircolors -b)"
zstyle ':completion:*:default' list-colors ${(s.:.)LS_COLORS}
zstyle ':completion:*' list-colors ''
zstyle ':completion:*' list-prompt %SAt %p: Hit TAB for more, or the character to insert%s
zstyle ':completion:*' matcher-list '' 'm:{a-z}={A-Z}' 'm:{a-zA-Z}={A-Za-z}' 'r:|[._-]=* r:|=* l:|=*'
zstyle ':completion:*' menu select=long
zstyle ':completion:*' select-prompt %SScrolling active: current selection at %p%s
zstyle ':completion:*' use-compctl false
zstyle ':completion:*' verbose true

zstyle ':completion:*:*:kill:*:processes' list-colors '=(#b) #([0-9]#)*=0=01;31'
zstyle ':completion:*:kill:*' command 'ps -u $USER -o pid,%cpu,tty,cputime,cmd'
# User configuration

export MANPATH="/usr/local/man:$MANPATH"

# You may need to manually set your language environment
# export LANG=en_US.UTF-8

# Preferred editor for local and remote sessions
export EDITOR='nvim'

# Compilation flags
# export ARCHFLAGS="-arch x86_64"

# Set personal aliases, overriding those provided by oh-my-zsh libs,
# plugins, and themes. Aliases can be placed here, though oh-my-zsh
# users are encouraged to define aliases within the ZSH_CUSTOM folder.
# For a full list of active aliases, run `alias`.

export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

# Generated for envman. Do not edit.
[ -s "$HOME/.config/envman/load.sh" ] && source "$HOME/.config/envman/load.sh"

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

# Autocomplete sources
complete -C /home/linuxbrew/.linuxbrew/Cellar/terraform/0.13.4/bin/terraform terraform
complete -o nospace -C /usr/bin/terraform terraform
autoload -U +X bashcompinit && bashcompinit
source <(kubectl completion zsh)
source /usr/share/doc/fzf/examples/completion.zsh

source /usr/share/doc/fzf/examples/key-bindings.zsh

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # set DISPLAY variable to the IP automatically assigned to WSL2
  export DISPLAY=$(route.exe print | grep 0.0.0.0 | head -1 | awk '{print $4}'):0.0

  # Used for vagrant - Enables vagrant use from within WSL2
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
  export PATH="$PATH:/mnt/c/Program Files/Oracle/VirtualBox"
fi
