#!/bin/bash
set -euo pipefail

# Install zsh itself (previously done by prerequisites.sh)
sudo apt-get update
sudo apt-get install -y zsh

if [[ ! -d "$HOME"/.local/share/zinit/zinit.git ]]; then
  bash -c "$(curl --fail --show-error --silent --location https://raw.githubusercontent.com/zdharma-continuum/zinit/HEAD/scripts/install.sh)"
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

# Enable italics and 256color for terminal
tic "$DIR"/xterm-256color-italic.terminfo

sudo apt-get install -y fonts-powerline powerline

# Switch default shell to zsh (previously done by bootstrap.sh)
if command -v zsh >/dev/null 2>&1 && command -v chsh >/dev/null 2>&1; then
  ZSH="$(command -v zsh)"
  if grep -qx "$ZSH" /etc/shells 2>/dev/null || true; then
    chsh -s "$ZSH" || true
  fi
fi

# Preserve any existing .zshrc before stow links ours
if [[ -f "$HOME/.zshrc" && ! -L "$HOME/.zshrc" ]]; then
  mv "$HOME/.zshrc" "$HOME/.zshrc_original"
fi
