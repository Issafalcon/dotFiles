#!/bin/bash

sudo apt-get update -y &&
	sudo apt-get install -y \
		texlive-full \
		latexmk \
		xdotool \
		xindy

# Need python and pip to install below
if [[ ! -v python3 ]]; then
	echo "Python 3 found. Skipping python 3 installation"
else
	"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "python"
	path+=(/usr/bin/pip3)
fi

pip install pygments
