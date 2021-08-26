#!/bin/zsh

brew install fzf

# To install useful key bindings and fuzzy completion:
$(brew --prefix)/opt/fzf/install

# Uses zinit fzf package and adds the man pages from fzf
# Also installs fzf-tmux
zinit pack=bgn atclone="cp fzy.1 $ZPFX/man/man1" for fzy

# Source the fzf script to get the zsh keybindings and config
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh
