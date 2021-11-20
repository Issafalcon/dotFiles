#!/bin/bash

ARGS=$(getopt -a --options d:rp --long "dev-proxy:,run,package" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/gp-dashboard/gp-dashboard-app -d -s gp-dashboard -n gp-dash)

tmux send-keys -t gp-dashboard:1.0 nvim C-m

tmux new-window -c "${PROJECTS}"/PatientDemandApp_API -t gp-dashboard:2 -n API

tmux send-keys -t gp-dashboard:2.0 nvim C-m

tmux new-window -c "${PROJECTS}"/PatientDemandApp_API/PatientDemandApp.WebAPI -t gp-dashboard:3 -n Terminals
tmux splitw -c "${PROJECTS}"/gp-dashboard/gp-dashboard-app -t gp-dashboard:3

tmux select-layout -t gp-dashboard:1 tiled
tmux select-pane -t gp-dashboard:1.0

tmux select-layout -t gp-dashboard:2 tiled
tmux select-pane -t gp-dashboard:2.0

tmux select-layout -t gp-dashboard:3 tiled
tmux select-pane -t gp-dashboard:3.0

while true; do
  case "$1" in
    -r|--run)
      tmux send-keys -t gp-dashboard:3.0 dotnet\ run C-m
      tmux send-keys -t gp-dashboard:3.1 npm\ run\ start C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t gp-dashboard
else
  tmux -u switch-client -t gp-dashboard
fi
