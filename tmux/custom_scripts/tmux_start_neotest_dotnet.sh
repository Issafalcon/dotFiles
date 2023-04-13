#!/bin/bash

# Check if wiki dir exists and create it if not.
if [ ! -d "${PROJECTS}"/neotest-dotnet ]; then

  # Clone repo into dir.
  git clone https://github.com/Issafalcon/neotest-dotnet.git "${PROJECTS}"/neotest-dotnet
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/neotest-dotnet -d -s neotest-dotnet -n neotest-dotnet)

tmux send-keys -t neotest-dotnet:1.0 nvim C-m

tmux select-layout -t neotest-dotnet:1 tiled
tmux select-pane -t neotest-dotnet:1.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t neotest-dotnet
else
  tmux -u switch-client -t neotest-dotnet
fi
