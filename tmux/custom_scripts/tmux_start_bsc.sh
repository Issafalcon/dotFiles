#!/bin/bash

ARGS=$(getopt -a --options d:r --long "dev-proxy:,run" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/BSC/Frontend -d -s BSC -n Frontend)

tmux send-keys -t BSC:Frontend.0 nvim C-m

# Create other windows.
tmux new-window -c "${PROJECTS}"/BSC/Backend/BSC.Api -t BSC:2 -n API
tmux new-window -c "${PROJECTS}"/BSC/Frontend -t BSC:3 -n Terminals

tmux send-keys -t BSC:API.0 nvim C-m

tmux select-layout -t BSC:Frontend tiled
tmux select-pane -t BSC:Frontend.0

tmux select-layout -t BSC:API tiled
tmux select-pane -t BSC:API.0

tmux splitw -c "${PROJECTS}"/BSC/Backend/BSC.Api -t BSC:Terminals
tmux select-layout -t BSC:Terminals tiled
tmux select-pane -t BSC:Terminals.0

while true; do
  case "$1" in
    -r|--run)
      # TODO: Run correct commands here
      tmux send-keys -t BSC:3.0 dotnet\ run C-m
      tmux send-keys -t BSC:3.1 yarn\ start C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t BSC
else
  tmux -u switch-client -t BSC
fi
