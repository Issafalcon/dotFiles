#!/bin/bash
set -euo pipefail

BASH=$(which bash) && command -v chsh >/dev/null 2>&1 && chsh -s "$BASH"
sudo apt-get remove -y fonts-powerline powerline zsh || true
rm -rf "$HOME/.local/share/zinit"
