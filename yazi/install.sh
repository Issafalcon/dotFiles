#!/bin/bash

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

# Install yazi and supporting previewer tools
brew install yazi \
  ImageMagick \
  ffmpeg \
  fd

git clone https://github.com/MasouShizuka/projects.yazi.git ~/.config/yazi/plugins/projects.yazi
