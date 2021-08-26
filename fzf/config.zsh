#!/bin/zsh

# Uses zinit fzf package and adds the man pages from fzf
# Also installs fzf-tmux
zinit pack=bgn atclone="cp fzy.1 $ZPFX/man/man1" for fzy

# Source the fzf script to get the zsh keybindings and config
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh
