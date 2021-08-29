#!/bin/bash

if [[ ! -d "$HOME"/.zinit ]]; then
	mkdir ~/.zinit
	git clone https://github.com/zdharma/zinit.git ~/.zinit/bin
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

# Enable italics and 256color for terminal
tic "$DIR"/xterm-256color-italic.terminfo

sudo apt-get update \
	&& sudo apt-get install fonts-powerline \
	&& sudo apt-get install powerline
