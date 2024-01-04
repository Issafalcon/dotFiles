#!/bin/bash

# Check if wiki dir exists and create it if not.
if [ ! -d "${PROJECTS}"/wiki ]; then

  # Clone wiki repo into wiki dir.
  git clone https://github.com/Issafalcon/wiki.git "${PROJECTS}"/wiki
fi

if [ ! -d "${PROJECTS}"/wiki-md ]; then

  # Clone wiki repo into wiki dir.
  git clone https://github.com/Issafalcon/wiki-md.git "${PROJECTS}"/wiki-md
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/wiki-md -d -s wiki -n Obsidian)

tmux send-keys -t wiki:1.0 nvim C-m

# Terminal
tmux new-window -c "${PROJECTS}"/wiki -t wiki:2 -n Terminals

# Old wiki
tmux new-window -c "${DOTFILES}"/wiki -t wiki:3 -n wiki
tmux send-keys -t config:3.0 nvim C-m

tmux select-layout -t wiki:1 tiled
tmux select-layout -t wiki:2 tiled

tmux select-window -t wiki:Obsidian
tmux select-pane -t wiki:Obsidian.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t wiki
else
  tmux -u switch-client -t wiki
fi
