#!/bin/bash

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/wcc-deployment -d -s SRE -n Deployments)

tmux send-keys -t SRE:1.0 nvim C-m

tmux new-window -c "${PROJECTS}"/wcc-env-config -t SRE:2 -n Env-Config

tmux send-keys -t SRE:2.0 nvim C-m

tmux new-window -c "${PROJECTS}"/wcc-env-config -t SRE:3 -n Terminals
tmux splitw -c "${PROJECTS}"/wcc-env-config -t SRE:3

tmux select-layout -t SRE:1 tiled
tmux select-pane -t SRE:1.0

tmux select-layout -t SRE:2 tiled
tmux select-pane -t SRE:2.0

tmux select-layout -t SRE:3 tiled
tmux select-pane -t SRE:3.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t SRE
else
  tmux -u switch-client -t SRE
fi
