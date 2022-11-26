#!/bin/bash

# Install supporting tools for Neovim and neovim plugins
sudo apt-get update
sudo apt-get install ripgrep
sudo apt-get install cmake
sudo apt-get install automake
sudo apt-get install ninja-build
sudo apt-get install silversearcher-ag # Install silversearcher to perform fzf searches using ag insteak of ack
sudo apt-get install -y exuberant-ctags
sudo apt-get install clang
sudo apt-get install sqlite3 libsqlite3-dev # For persisting history of yanks between sessions

# Needed for ueberzug
sudo apt-get install libjpeg8-dev
sudo apt-get install zlib1g-dev
sudo apt-get install libxtst-dev
sudo apt-get install libx11
sudo apt-get install libxext-dev
sudo apt-get install xllproto-xext-dev
sudo apt-get install libtool-bin
sudo apt-get install gettext

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
sudo apt-get install chktex

# Install debug adapters: Use Mason plugin to install others
mkdir -p ~/.local/share/nvim/mason/packages

git clone https://github.com/rogalmic/vscode-bash-debug.git ~/.local/share/nvim/mason/packages
cd ~/.local/share/nvim/mason/packages/vscode-bash-debug
npm install
npm run compile
