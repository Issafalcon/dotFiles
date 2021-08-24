#!/bin/bash

if [[ ! -d "$HOME"/.zinit ]]; then
	mkdir ~/.zinit
	git clone https://github.com/zdharma/zinit.git ~/.zinit/bin
fi

sudo apt-get update \
	&& sudo apt-get install fonts-powerline \
	&& sudo apt-get install powerline

