#!/bin/bash

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c ~/dotFiles -d -s config -n dotfiles)

tmux send-keys -t config:1.0 nvim C-m

# Neovim Config
tmux new-window -c "${DOTFILES}"/nvim/.config/nvim -t config:2 -n nvim-config
tmux send-keys -t config:2.0 nvim C-m

# Teminals Window
tmux new-window -c ~/dotFiles -t config:3 -n Terminals
tmux splitw -c "${DOTFILES}"/nvim/.config/nvim -t config:3 -l 50%

tmux select-layout -t config:1 tiled
tmux select-layout -t config:2 tiled

tmux select-window -t config:nvim-config
tmux select-pane -t config:nvim-config.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t config
else
  tmux -u switch-client -t config
fi
