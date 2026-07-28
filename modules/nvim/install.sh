#!/bin/bash

# Install supporting tools for Neovim and neovim plugins
sudo apt-get update
sudo apt-get install -y ripgrep \
  cmake \
  automake \
  ninja-build \
  silversearcher-ag \
  exuberant-ctags \
  clang \
  sqlite3 libsqlite3-dev

sudo apt-get install -y libjpeg8-dev \
  zlib1g-dev \
  libxtst-dev \
  libx11 \
  libxext-dev \
  xllproto-xext-dev \
  libtool-bin \
  gettext \
  lua5.1 \
  liblua5.1-dev

# Set python virtual env
if [[ ! -d "$HOME/python3/envs/neovim" ]]; then
  mkdir -p "$HOME"/python3/envs
  cd "$HOME"/python3/envs || exit
  python3 -m venv neovim
  source "${HOME}"/python3/envs/neovim/bin/activate
  python3 -m pip install pynvim
  python3 -m pip install neovim
  python3 -m pip install neovim-remote
  python3 -m pip install jupyter ipykernel nbclient nbformat jupyter-cache PyYAML matplotlib plotly
  python3 -m pip install pandas
  deactivate
fi

# Register the neovim venv's ipykernel as a user-wide Jupyter kernel so that
# Molten can launch it without the venv being active
mkdir -p "$HOME/.local/share/jupyter/runtime"
source "${HOME}"/python3/envs/neovim/bin/activate
python3 -m ipykernel install --user --name python3 --display-name "Python 3 (neovim)"
deactivate

# pip3 no longer can install global packages
# Below is needed for rnvimr
sudo apt install python3-pynvim
npm install -g tree-sitter-cli --allow-scripts=tree-sitter-cli
npm install -g neovim

# Get Neovim latest release as app image and move to /usr/bin/nvim
sudo curl -Lo /usr/bin/nvim https://github.com/neovim/neovim/releases/download/v0.12.0/nvim-linux-x86_64.appimage
sudo chmod 777 /usr/bin/nvim

# Install formatters / linters for LSP
sudo apt-get install -y chktex

# Install debuggers that require manual install

## vscode-js-debug
mkdir -p "$HOME"/.local/share/nvim
cd "$HOME"/.local/share/nvim || exit
git clone https://github.com/microsoft/vscode-js-debug
cd vscode-js-debug || exit
npm install --legacy-peer-deps
npx gulp vsDebugServerBundle
mv dist out

# FUSE needed to run app images
sudo apt install -y libfuse2

# Install ast-grep for grug-far search and replace
brew install ast-grep
