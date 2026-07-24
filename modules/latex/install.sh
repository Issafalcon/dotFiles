#!/bin/bash

sudo apt-get update -y &&
	sudo apt-get install -y \
		texlive-full \
		latexmk \
		xdotool \
		xindy

# Need python and pip to install below
if command -v python3 >/dev/null; then
	echo "Python 3 found. Skipping python 3 installation"
else
	echo "python module is required. Install it from the TUI first." >&2
	exit 1
fi

pip install pygments
