#!/bin/zsh
# uncomment this and the last line for zprof info
# zmodload zsh/zprof

# +------------------+
# | PACKAGE MANAGER  |
# +------------------+

source "$HOME/.zinit/bin/zinit.zsh"
autoload -Uz _zinit
(( ${+_comps} )) && _comps[zinit]=_zinit

# Load a few important annexes, without Turbo
# (this is currently required for annexes)
zinit light-mode for \
    zinit-zsh/z-a-rust \
    zinit-zsh/z-a-as-monitor \
    zinit-zsh/z-a-patch-dl \
    zinit-zsh/z-a-bin-gem-node

### Added by Zinit's installer
if [[ ! -f $HOME/.zinit/bin/zinit.zsh ]]; then
    print -P "%F{33}▓▒░ %F{220}Installing %F{33}DHARMA%F{220} Initiative Plugin Manager (%F{33}zdharma/zinit%F{220})…%f"
    command mkdir -p "$HOME/.zinit" && command chmod g-rwX "$HOME/.zinit"
    command git clone https://github.com/zdharma/zinit "$HOME/.zinit/bin" && \
        print -P "%F{33}▓▒░ %F{34}Installation successful.%f%b" || \
        print -P "%F{160}▓▒░ The clone has failed.%f%b"
fi

### End of Zinit's installer chunk

# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi


# +------------+
# | FUNCTIONS  |
# +------------+

# allow functions to have local options
setopt LOCAL_OPTIONS
# allow functions to have local traps
setopt LOCAL_TRAPS

fpath=("$DOTFILES/zsh/functions" "${fpath[@]}")
autoload -U "${DOTFILES}"/zsh/functions/*(.:t)

# +---------+
# | GENERAL |
# +---------+

# don't nice background tasks
setopt NO_BG_NICE
setopt NO_HUP
setopt NO_BEEP
setopt promptsubst

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

# +-----------------------------+
# | COMPLETION AND HIGHLIGHTING |
# +-----------------------------+

# zstyle pattern for the completion
# :completion:<function>:<completer>:<command>:<argument>:<tag>

# Should be called before compinit
zmodload zsh/complist

zinit wait lucid for \
  atinit"zicompinit; zicdreplay" \
      zdharma/fast-syntax-highlighting \
  atload"_zsh_autosuggest_start" \
      zsh-users/zsh-autosuggestions \
  blockf atpull'zinit creinstall -q .' \
      zsh-users/zsh-completions

# Use hjlk in menu selection (during completion)
# Doesn't work well with interactive mode
bindkey -M menuselect 'h' vi-backward-char
bindkey -M menuselect 'k' vi-up-line-or-history
bindkey -M menuselect 'j' vi-down-line-or-history
bindkey -M menuselect 'l' vi-forward-char

bindkey -M menuselect '^xg' clear-screen
bindkey -M menuselect '^xi' vi-insert                      # Insert
bindkey -M menuselect '^xh' accept-and-hold                # Hold
bindkey -M menuselect '^xn' accept-and-infer-next-history  # Next
bindkey -M menuselect '^xu' undo                           # Undo

_comp_options+=(globdots) # With hidden files

source "$DOTFILES"/zsh/completion.zsh

# Only work with the Zsh function vman
# See $DOTFILES/zsh/scripts.zsh
# compdef vman="man"

# +-------+
# | THEME |
# +-------+

zinit ice depth=1; zinit light romkatv/powerlevel10k

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

# +---------+
# | MODULES |
# +---------+

# Check what modules have been 'installed'
MODULES=($(cat $HOME/.dotFileModules))
# Set and autoload all custom module files
for module in ${MODULES}; do
  [ -d "$DOTFILES/$module/functions" ] \
    && fpath=("$DOTFILES/$module/functions" "${fpath[@]}") \
    && autoload -U "${DOTFILES}"/"${module}"/functions/*(.:t)

  [ -f "$DOTFILES/$module/path.zsh" ] \
    && source "$DOTFILES/$module/path.zsh"

  # Load the zsh script snippets first
  [ -f "$DOTFILES/$module/scripts.zsh" ] \
    && source "$DOTFILES/$module/scripts.zsh"

  # Then add the custom_scripts to the path
  [ -d "$DOTFILES/$module/custom_scripts" ] \
    && path+=("${DOTFILES}"/"${module}"/custom_scripts)

  # Add specific zsh configuration
  [ -f "$DOTFILES/$module/config.zsh" ] \
    && source "$DOTFILES/$module/config.zsh"
done

# Adds homebrew path
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
