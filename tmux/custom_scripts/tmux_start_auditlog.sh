#!/bin/bash

ARGS=$(getopt -a --options d:r --long "dev-proxy:,run" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -d -s AuditLog -n ClientApp)

# Create other windows.
tmux new-window -c ~/repos/wcc-auditlog/Waters.Cloud.AuditLog -t AuditLog:2 -n Api
tmux new-window -c ~/repos/wcc-auditlog -t AuditLog:3 -n Terminals

tmux select-layout -t AuditLog:1 tiled
tmux select-layout -t AuditLog:1
tmux select-pane -t AuditLog:1.0

tmux select-layout -t AuditLog:2 tiled
tmux select-layout -t AuditLog:2
tmux select-pane -t AuditLog:2.0

tmux splitw -c ~/repos/wcc-auditlog -t AuditLog:3
tmux select-layout -t AuditLog:3 tiled
tmux select-layout -t AuditLog:3 tiled
tmux select-layout -t AuditLog:3 tiled
tmux select-pane -t AuditLog:3.0

tmux send-keys -t AuditLog:1 cd\ ~/repos/wcc-auditlog C-m
tmux send-keys -t AuditLog:1.0 cd\ AuditLog/ClientApp C-m
tmux send-keys -t AuditLog:1.0 nvim C-m

tmux send-keys -t AuditLog:2.0 nvim C-m

tmux send-keys -t AuditLog:3.0 cd\ Waters.Cloud/AuditLog C-m
tmux send-keys -t AuditLog:3.1 cd\ Waters.Cloud.AuditLog/ClientApp C-m

while true; do
  case "$1" in
    -d|--dev-proxy)
      tmux new-window -c ~/repos/wcc-auditlog -t AuditLog:4 -n Dev-Proxy
      tmux send-keys -t AuditLog:4 dev-proxy\ -l\ "${2}" C-m
      shift 2;;
    -r|--run)
      tmux send-keys -t AuditLog:3.0 dotnet\ run C-m
      tmux send-keys -t AuditLog:3.1 ng\ serve C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t AuditLog
else
  tmux -u switch-client -t AuditLog
fi

