#!/bin/bash

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  echo "homebrew module is required. Install it from the TUI first." >&2
  exit 1
fi

# Install ImageMagick via its own module
if command -v magick >/dev/null; then
  echo "ImageMagick found. Skipping imagemagick installation"
else
  echo "imagemagick module is required. Install it from the TUI first." >&2
  exit 1
fi

# Install yazi and supporting previewer tools
brew install yazi \
  ffmpeg \
  fd

sudo git clone https://github.com/MasouShizuka/projects.yazi.git "${SCRIPT_DIR}"/.config/yazi/plugins/projects.yazi
