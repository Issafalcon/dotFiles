#!/bin/bash
set -euo pipefail

sudo rm -f /usr/bin/nvim
rm -rf "$HOME/.local/share/nvim/vscode-js-debug"
rm -rf "$HOME/python3/envs/neovim"
