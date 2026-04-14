#!/bin/bash

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

# Install ImageMagick via its own module
if command -v magick >/dev/null; then
  echo "ImageMagick found. Skipping imagemagick installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "imagemagick"
fi

# Install yazi and supporting previewer tools
brew install yazi \
  ffmpeg \
  fd

sudo git clone https://github.com/MasouShizuka/projects.yazi.git "${SCRIPT_DIR}"/.config/yazi/plugins/projects.yazi
