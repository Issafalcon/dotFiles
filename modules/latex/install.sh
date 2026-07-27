#!/bin/bash

sudo apt-get update -y &&
	sudo apt-get install -y \
		texlive-full \
		latexmk \
		xdotool \
		xindy

pip install pygments
