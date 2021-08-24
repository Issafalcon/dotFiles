#!/bin/bash

if [[ ! -d "$HOME"/.zinit ]]; then
	mkdir ~/.zinit
	git clone https://github.com/zdharma/zinit.git ~/.zinit/bin
fi
