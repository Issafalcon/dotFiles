#!/bin/bash

# Install supporting tools for Neovim and neovim plugins
sudo apt-get update
sudo apt-get install -y ripgrep
 cmake
 automake
 ninja-build
 silversearcher-ag # Install silversearcher to perform fzf searches using ag insteak of ack
 exuberant-ctags
 clang
 sqlite3 libsqlite3-dev # For persisting history of yanks between sessions

# Needed for ueberzug
sudo apt-get install -y libjpeg8-dev
 zlib1g-dev
 libxtst-dev
 libx11
 libxext-dev
 xllproto-xext-dev
 libtool-bin
 gettext

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Need python and pip to install below
if [[ ! -v python3 ]]; then
	echo "Python 3 found. Skipping python 3 installation"
else
	"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "python"
	path+=(/usr/bin/pip3)
fi

# Also need to use node for npm
if [[ ! -v node ]]; then
	echo "Node found. Skipping python 3 installation"
else
	"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "node"
fi

# Install go
if [[ ! -v go ]]; then
	echo "go found. Skipping go installation"
else
	"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "go"
fi

pip3 install --user neovim-remote
pip3 install --user ueberzug
pip3 install --user pynvim
npm install -g tree-sitter-cli

# Get Neovim latest release as app image and move to /usr/bin/nvim
sudo curl -Lo /usr/bin/nvim https://github.com/neovim/neovim/releases/download/v0.8.0/nvim.appimage
sudo chmod 777 /usr/bin/nvim

# Install formatters / linters for LSP
sudo apt-get install -y chktex

# FUSE needed to run app images
sudo apt install -y fuse libfuse2