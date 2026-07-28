#!/bin/bash

# Check if scrimmage dir exists and create it if not.
if [ ! -d "${PROJECTS}"/scrimmage-api ]; then
  # Clone scrimmage repo into scrimmage dir.
  git clone https://github.com/Issafalcon/scrimmage-api.git "${PROJECTS}"/scrimmage-api
fi

if [ ! -d "${PROJECTS}"/scrimmage-web ]; then
  # Clone scrimmage repo into scrimmage dir.
  git clone https://github.com/Issafalcon/scrimmage-web.git "${PROJECTS}"/scrimmage-web
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/scrimmage-api -d -s scrimmage -n API)

tmux send-keys -t scrimmage:1.0 nvim C-m

# Terminal
tmux new-window -c "${PROJECTS}"/scrimmage-web -t scrimmage:2 -n Web

tmux send-keys -t scrimmage:Web.1 nvim C-m

tmux new-window -c "${PROJECTS}"/scrimmage-api -t scrimmage:3 -n Terminals
tmux split-window -c "${PROJECTS}"/scrimmage-web -t scrimmage:Terminals -h -l 30%

tmux select-window -t scrimmage:Web
tmux select-pane -t scrimmage:Web.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t scrimmage
else
  tmux -u switch-client -t scrimmage
fi
