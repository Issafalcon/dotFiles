#!/bin/bash

ARGS=$(getopt -a --options d:rp --long "dev-proxy:,run,package" -- "$@")

eval set -- "$ARGS"

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -d -s Notifications-MFE -n Center)

# Create other windows.
tmux new-window -c ~/repos/wcc-mfe-notifications/notifications-mfe -t Notifications-MFE:2 -n MFE
tmux new-window -c ~/repos/wcc-mfe-notifications -t Notifications-MFE:3 -n Terminals

tmux select-layout -t Notifications-MFE:1 tiled
tmux select-layout -t Notifications-MFE:1
tmux select-pane -t Notifications-MFE:1.0

tmux select-layout -t Notifications-MFE:2 tiled
tmux select-layout -t Notifications-MFE:2
tmux select-pane -t Notifications-MFE:2.0

tmux splitw -c ~/repos/wcc-mfe-notifications -t Notifications-MFE:3
tmux splitw -h -c ~/repos/wcc-mfe-notifications -t Notifications-MFE:3.1
tmux select-layout -t Notifications-MFE:3 tiled
tmux select-pane -t Notifications-MFE:3.0

sleep 3

tmux send-keys -t Notifications-MFE:1 cd\ ~/repos/wcc-mfe-notifications/notifications-center C-m
tmux send-keys -t Notifications-MFE:1.0 nvim C-m

sleep 1

tmux send-keys -t Notifications-MFE:2.0 nvim C-m

sleep 1

tmux send-keys -t Notifications-MFE:3.0 cd\ notifications-center/Waters.Cloud.NotificationsCenter C-m
tmux send-keys -t Notifications-MFE:3.1 cd\ notifications-center/Waters.Cloud.NotificationsCenter/ClientApp C-m

while true; do
  case "$1" in
    -d|--dev-proxy)
      tmux new-window -c ~/repos/wcc-mfe-notifications -t Notifications-MFE:4 -n Dev-Proxy
      tmux send-keys -t Notifications-MFE:4 dev-proxy\ -l\ "${2}" C-m
      shift 2;;
    -r|--run)
      tmux send-keys -t Notifications-MFE:3.0 dotnet\ run C-m
      tmux send-keys -t Notifications-MFE:3.1 ng\ serve C-m
      shift;;
    -p|--package)
      tmux send-keys -t Notifications-MFE:3.2 npm\ run\ build-package:all-mfe C-m
      shift;;
    --)
      break;;
  esac
done

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t Notifications-MFE
else
  tmux -u switch-client -t Notifications-MFE
fi

