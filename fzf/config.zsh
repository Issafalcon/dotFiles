#!/bin/zsh

# Uses zinit fzf package and adds the man pages from fzf
# Also installs fzf-tmux
zinit pack=bgn atclone="cp fzy.1 $ZPFX/man/man1" for fzy

source /usr/share/doc/fzf/examples/key-bindings.zsh
source /usr/share/doc/fzf/examples/completion.zsh
