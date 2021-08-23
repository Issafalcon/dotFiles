#!/bin/bash

ARGS=$(getopt -a --options d:r --long "dev-proxy:,run" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c ~/repos/wcc-tenancyas -d -s TenancyAS -n tenancyas-root)

sleep 4

tmux send-keys -t TenancyAS:1.0 nvim C-m

# Create other windows.
tmux new-window -c ~/repos/wcc-tenancyas/Services/Waters.TenancyAS.WebAPI -t TenancyAS:2 -n Api
tmux new-window -c ~/repos/wcc-tenancyas/Services/Waters.TenancyAS.WebAPI -t TenancyAS:3 -n Terminals

sleep 3

tmux send-keys -t TenancyAS:2.0 nvim C-m

tmux select-layout -t TenancyAS:1 tiled
tmux select-pane -t TenancyAS:1.0

tmux select-layout -t TenancyAS:2 tiled
tmux select-pane -t TenancyAS:2.0

tmux splitw -c ~/repos/wcc-tenancyas -t TenancyAS:3
tmux select-layout -t TenancyAS:3 tiled
tmux select-pane -t TenancyAS:3.0

while true; do
  case "$1" in
    -d|--dev-proxy)
      tmux new-window -c ~/repos/wcc-tenancyas -t TenancyAS:4 -n Dev-Proxy

      sleep 3

      tmux send-keys -t TenancyAS:4 dev-proxy\ -l\ "${2}" C-m
      shift 2;;
    -r|--run)

      sleep 3

      tmux send-keys -t TenancyAS:3.0 dotnet\ run C-m
      tmux send-keys -t TenancyAS:3.1 ng\ serve C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t TenancyAS
else
  tmux -u switch-client -t TenancyAS
fi

