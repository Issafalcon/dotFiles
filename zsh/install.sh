#!/bin/bash

if [[ ! -d "$HOME"/.local/share/zinit/zinit.git ]]; then
  bash -c "$(curl --fail --show-error --silent --location https://raw.githubusercontent.com/zdharma-continuum/zinit/HEAD/scripts/install.sh)"
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

# Enable italics and 256color for terminal
tic "$DIR"/xterm-256color-italic.terminfo

sudo apt-get update \
	&& sudo apt-get install fonts-powerline \
	&& sudo apt-get install powerline
