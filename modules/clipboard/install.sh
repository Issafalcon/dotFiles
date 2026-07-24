#!/bin/bash

# Install win yank if on WSL so vim can use windows clipboard
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  echo "Not installing clipboard on WSL as win32yank is part of Neovim"
else
  sudo apt-get install -y xclip
fi
