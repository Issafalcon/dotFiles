#!/bin/bash

ARGS=$(getopt -a --options d:r --long "dev-proxy:,run" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c ~/repos/wcc-cloudhub/Waters.Cloud.Hub/ClientApp -d -s CloudHub -n ClientApp)

# Create other windows.
tmux new-window -c ~/repos/wcc-cloudhub/Waters.Cloud.Hub -t CloudHub:2 -n Api
tmux new-window -c ~/repos/wcc-cloudhub -t CloudHub:3 -n Terminals

tmux select-layout -t CloudHub:1 tiled
tmux select-pane -t CloudHub:1.0

tmux select-layout -t CloudHub:2 tiled
tmux select-pane -t CloudHub:2.0

tmux splitw -c ~/repos/wcc-cloudhub -t CloudHub:3
tmux select-layout -t CloudHub:3 tiled
tmux select-pane -t CloudHub:3.0

tmux send-keys -t CloudHub:1.0 nvim C-m
tmux send-keys -t CloudHub:2.0 nvim C-m
tmux send-keys -t CloudHub:3.0 cd\ Waters.Cloud.Hub C-m
tmux send-keys -t CloudHub:3.1 cd\ Waters.Cloud.Hub/ClientApp C-m

while true; do
  case "$1" in
    -d|--dev-proxy)
      tmux new-window -c ~/repos/wcc-cloudhub -t CloudHub:4 -n Dev-Proxy
      
      tmux send-keys -t CloudHub:4 dev-proxy\ -l\ "${2}" C-m
      shift 2;;
    -r|--run)
      tmux send-keys -t CloudHub:3.0 dotnet\ run C-m
      tmux send-keys -t CloudHub:3.1 ng\ serve C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t CloudHub
else
  tmux -u switch-client -t CloudHub
fi


