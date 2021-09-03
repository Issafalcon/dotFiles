#!/bin/bash

ARGS=$(getopt -a --options d:r --long "dev-proxy:,run" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/wcc-tenancyadmin/TenancyAdmin/ClientApp -d -s TenancyAdmin -n ClientApp)

tmux send-keys -t TenancyAdmin:1.0 nvim C-m

# Create other windows.
tmux new-window -c "${PROJECTS}"/wcc-tenancyadmin/TenancyAdmin -t TenancyAdmin:2 -n Api
tmux new-window -c "${PROJECTS}"/wcc-tenancyadmin/TenancyAdmin -t TenancyAdmin:3 -n Terminals

tmux send-keys -t TenancyAdmin:2.0 nvim C-m

tmux select-layout -t TenancyAdmin:1 tiled
tmux select-pane -t TenancyAdmin:1.0

tmux select-layout -t TenancyAdmin:2 tiled
tmux select-pane -t TenancyAdmin:2.0

tmux splitw -c "${PROJECTS}"/wcc-tenancyadmin/TenancyAdmin/ClientApp -t TenancyAdmin:3
tmux select-layout -t TenancyAdmin:3 tiled
tmux select-pane -t TenancyAdmin:3.0

while true; do
  case "$1" in
    -d|--dev-proxy)
      tmux new-window -c "${PROJECTS}"/wcc-tenancyadmin -t TenancyAdmin:4 -n Dev-Proxy
      tmux send-keys -t TenancyAdmin:4 dev-proxy\ -l\ "${2}" C-m
      shift 2;;
    -r|--run)
      tmux send-keys -t TenancyAdmin:3.0 dotnet\ run C-m
      tmux send-keys -t TenancyAdmin:3.1 ng\ serve C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t TenancyAdmin
else
  tmux -u switch-client -t TenancyAdmin
fi

