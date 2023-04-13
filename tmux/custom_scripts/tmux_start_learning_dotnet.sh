#!/bin/bash

# Check if wiki dir exists and create it if not.
if [ ! -d "${PROJECTS}"/learning-dotnet ]; then

  # Clone repo into dir.
  git clone https://github.com/Issafalcon/learning-dotnet.git "${PROJECTS}"/learning-dotnet
fi

# Create the session and the first window. Manually switch to root
# directory if required to support tmux < 1.9
TMUX=$(tmux new-session -c "${PROJECTS}"/learning-dotnet/UnitTesting -d -s learning-dotnet -n unit-testing)

tmux send-keys -t learning-dotnet:1.0 nvim C-m

tmux select-layout -t learning-dotnet:1 tiled
tmux select-pane -t learning-dotnet:1.0

if [ -z "$TMUX" ]; then
  tmux -u attach-session -t learning-dotnet
else
  tmux -u switch-client -t learning-dotnet
fi
