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

# Needed for ueberzug
sudo apt-get install libjpeg8-dev 
sudo apt-get install zlib1g-dev 
sudo apt-get install libxtst-dev 
sudo apt-get install libx11 
sudo apt-get install libxext-dev 
sudo apt-get install xllproto-xext-dev 
sudo apt-get install libtool-bin 
sudo apt-get install gettext 

# Ranger plugins for dev-icons
git clone https://github.com/alexanderjeurissen/ranger_devicons ~/.config/ranger/plugins/ranger_devicons

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Need python and pip to install below
"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "python"

# Also need to use node for npm
"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "node"

# Install homebrew
"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"

pip3 install --user neovim-remote
pip3 install --user ueberzug
pip3 install --user pynvim
npm install -g tree-sitter-cli

# Get Neovim latest release as app image and move to /usr/bin/nvim
sudo curl -Lo /usr/bin/nvim https://github.com/neovim/neovim/releases/latest/download/nvim.appimage
sudo chmod 777 nvim

# Install language servers that can't be installed via the LspInstall vim command (via lspinstall plugin)
brew install efm-langserver
brew install texlab
brew install hashicorp/tap/terraform-ls
npm install -g emmet-ls
npm install -g stylelint-lsp

# Install formatters / linters for LSP
npm install -g lua-fmt
npm install -g eslint
npm install -g eslint_d
npm install -g prettier
npm install -g markdownlint
sudo apt-get install shellcheck
sudo apt install yamllint
brew install shfmt

# Install debug adapters - Used for DAP only. Vimspector installs them as 'gadgets'
# mkdir -p ~/debug-adapters
#
# git clone https://github.com/microsoft/vscode-node-debug2.git ~/debug-adapters/vscode-node-debug2
# cd ~/debug-adapters/vscode-node-debug2
# npm install
#
# git clone https://github.com/Microsoft/vscode-chrome-debug ~/debug-adapters/vscode-chrome-debug
# cd ~/debug-adapters/vscode-chrome-debug
# npm install
# npm run build
#
# git clone https://github.com/Samsung/netcoredbg.git ~/debug-adapters/netcoredbg
# cd ~/debug-adapters/netcoredbg
# mkdir build
# cd build
# CC=clang CXX=clang++ cmake ..
# make
# sudo make install
