#!/bin/bash

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c ~/repos/scripts -d -s scripts -n editor)

tmux send-keys -t scripts:1.0 nvim C-m

# Create other windows.
tmux new-window -c ~/repos/scripts -t scripts:2 -n Terminals

tmux select-layout -t scripts:1 tiled
tmux select-layout -t scripts:2 tiled

tmux select-pane -t scripts:2.0
tmux select-window -t scripts:editor
tmux select-pane -t scripts:editor.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t scripts
else
  tmux -u switch-client -t scripts
fi
