#!/bin/bash

if [ ! -d "${PROJECTS}"/neotest-dotnet ]; then
  # Clone repo into dir.
  git clone https://github.com/Issafalcon/neotest-dotnet.git "${PROJECTS}"/neotest-dotnet
fi

if [ ! -d "${PROJECTS}"/neotest ]; then
  # Clone repo into dir.
  git clone https://github.com/nvim-neotest/neotest.git "${PROJECTS}"/neotest
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/neotest-dotnet -d -s neotest-dotnet -n neotest-dotnet)

tmux send-keys -t neotest-dotnet:1.0 nvim C-m

# neotest core
tmux new-window -c "${PROJECTS}"/neotest -t neotest-dotnet:2 -n neotest
tmux send-keys -t neotest-dotnet:neotest nvim C-m

tmux select-window -t neotest-dotnet:neotest-dotnet
tmux select-pane -t neotest-dotnet:1.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t neotest-dotnet
else
  tmux -u switch-client -t neotest-dotnet
fi
