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

# Switch default login shell to zsh (previously done by bootstrap.sh).
# Must use sudo: streaming install has no TTY for chsh's PAM password prompt;
# the TUI already ran sudo -v before this script.
ZSH="$(command -v zsh)"
if [[ -n "$ZSH" ]]; then
  if ! grep -qx "$ZSH" /etc/shells 2>/dev/null; then
    echo "$ZSH" | sudo tee -a /etc/shells >/dev/null
  fi
  current="$(getent passwd "${USER:-$(id -un)}" | cut -d: -f7)"
  if [[ "$current" != "$ZSH" ]]; then
    sudo chsh -s "$ZSH" "${USER:-$(id -un)}"
    echo "set $($ZSH --version) at $ZSH as default shell"
  fi
fi

# Preserve any existing .zshrc before stow links ours
if [[ -f "$HOME/.zshrc" && ! -L "$HOME/.zshrc" ]]; then
  mv "$HOME/.zshrc" "$HOME/.zshrc_original"
fi
