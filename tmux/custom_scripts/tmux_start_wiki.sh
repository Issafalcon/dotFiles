#!/bin/bash

# Check if wiki dir exists and create it if not.
if [ ! -d "${PROJECTS}"/wiki ]; then

  # Clone wiki repo into wiki dir.
  git clone https://github.com/Issafalcon/wiki.git "${PROJECTS}"/wiki
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/wiki -d -s wiki -n editor)

tmux send-keys -t wiki:1.0 nvim C-m

# Create other windows.
tmux new-window -c "${PROJECTS}"/wiki -t wiki:2 -n Terminals

tmux select-layout -t wiki:1 tiled
tmux select-layout -t wiki:2 tiled

tmux select-window -t wiki:editor
tmux select-pane -t wiki:editor.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t wiki
else
  tmux -u switch-client -t wiki
fi
