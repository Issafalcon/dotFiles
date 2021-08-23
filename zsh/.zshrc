#!/bin/zsh
# uncomment this and the last line for zprof info
# zmodload zsh/zprof

# Check what modules have been 'installed'
MODULES = $(cat $HOME/.dotFileModules)

# +------------+
# | FUNCTIONS  |
# +------------+

# allow functions to have local options
setopt LOCAL_OPTIONS
# allow functions to have local traps
setopt LOCAL_TRAPS

# Set and autoload all custom module functions 
for module in "${MODULES}"; do
  source "$DOTFILES/$module/path.zsh"

  fpath=("$DOTFILES/$module/functions" "${fpath[@]}")
  autoload -Uz "${DOTFILES}"/"${module}"/functions/*(:t)
done

# +-------+
# | PATHS |
# +-------+

# Add paths for modules
for module in "${MODULES}"; do
  source "$DOTFILES/$module/path.zsh"
done

# +---------+
# | SCRIPTS |
# +---------+

for module in "${MODULES}"; do
  # Load the zsh script snippets first
  source "$DOTFILES/$module/scripts.zsh"

  # Then add the custom_scripts to the path
  path+=("$DOTFILES"/$module/custom_scripts)
done

# +---------+
# | GENERAL |
# +---------+

# don't nice background tasks
setopt NO_BG_NICE
setopt NO_HUP
setopt NO_BEEP

# +------------+
# | NAVIGATION |
# +------------+

setopt AUTO_CD              # Go to folder path without using cd.

setopt AUTO_PUSHD           # Push the old directory onto the stack on cd.
setopt PUSHD_IGNORE_DUPS    # Do not store duplicates in the stack.
setopt PUSHD_SILENT         # Do not print the directory stack after pushd or popd.

setopt CORRECT              # Spelling correction
setopt CDABLE_VARS          # Change directory to a path stored in a variable.
setopt EXTENDED_GLOB        # Use extended globbing syntax.

bd # Custom function in zsh module

# +---------+
# | HISTORY |
# +---------+

setopt EXTENDED_HISTORY          # Write the history file in the ':start:elapsed;command' format.
setopt SHARE_HISTORY             # Share history between all sessions.
setopt HIST_EXPIRE_DUPS_FIRST    # Expire a duplicate event first when trimming history.
setopt HIST_IGNORE_DUPS          # Do not record an event that was just recorded again.
setopt HIST_IGNORE_ALL_DUPS      # Delete an old recorded event if a new event is a duplicate.
setopt HIST_FIND_NO_DUPS         # Do not display a previously found event.
setopt HIST_IGNORE_SPACE         # Do not record an event starting with a space.
setopt HIST_SAVE_NO_DUPS         # Do not write a duplicate event to the history file.
setopt HIST_VERIFY               # Do not execute immediately upon history expansion.

autoload -U up-line-or-beginning-search
autoload -U down-line-or-beginning-search
autoload -U edit-command-line

#Configurations options

 zle -N up-line-or-beginning-search
 zle -N down-line-or-beginning-search
 zle -N edit-command-line

 # fuzzy find: start to type
 bindkey "$terminfo[kcuu1]" up-line-or-beginning-search
 bindkey "$terminfo[kcud1]" down-line-or-beginning-search
 bindkey "$terminfo[cuu1]" up-line-or-beginning-search
 bindkey "$terminfo[cud1]" down-line-or-beginning-search

# +----------+
# | VIM MODE |
# +----------+

# Vi mode
bindkey -v
export KEYTIMEOUT=1

# Change cursor
cursor_mode # Custom function in zsh module

# edit current command line with vim (vim-mode, then v)
autoload -Uz edit-command-line
zle -N edit-command-line
bindkey -M vicmd v edit-command-line

# +------------+
# | COMPLETION |
# +------------+

for module in "${MODULES}"; do
  source "$DOTFILES/$module/completion.zsh"
done

# +-----+
# | FZF |
# +-----+

 # search history with fzf if installed, default otherwise
 if test -d /usr/local/opt/fzf/shell; then
   # shellcheck disable=SC1091
   . /usr/local/opt/fzf/shell/key-bindings.zsh
 else
   bindkey '^R' history-incremental-search-backward
 fi

# +-----+
# | nvm |
# +-----+

[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

# Generated for envman. Do not edit.
[ -s "$HOME/.config/envman/load.sh" ] && source "$HOME/.config/envman/load.sh"


