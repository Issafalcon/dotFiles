#!/bin/bash

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
	echo "homebrew module is required. Install it from the TUI first." >&2
	exit 1
fi

# Install yazi and supporting previewer tools
brew install yazi \
  ImageMagick \
  ffmpeg \
  fd

sudo git clone https://github.com/MasouShizuka/projects.yazi.git "${SCRIPT_DIR}"/.config/yazi/plugins/projects.yazi
