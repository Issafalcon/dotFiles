#!/bin/bash

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Install yazi and supporting previewer tools
brew install yazi \
  ffmpeg \
  fd

sudo git clone https://github.com/MasouShizuka/projects.yazi.git "${SCRIPT_DIR}"/.config/yazi/plugins/projects.yazi
