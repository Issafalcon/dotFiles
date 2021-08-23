#!/bin/bash

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c ~/dotFiles -d -s dotfiles -n editor)

sleep 4

tmux send-keys -t dotfiles:1.0 nvim C-m

# Create other windows.
tmux new-window -c ~/dotFiles -t dotfiles:2 -n Terminals

tmux select-layout -t dotfiles:1 tiled
tmux select-layout -t dotfiles:2 tiled

tmux select-pane -t dotfiles:2.0
tmux select-window -t dotfiles:editor
tmux select-pane -t dotfiles:editor.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t dotfiles
else
  tmux -u switch-client -t dotfiles
fi
