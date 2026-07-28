#!/bin/bash
set -euo pipefail

# Restore bash as login shell (needs sudo — same reason as install.sh)
BASH="$(command -v bash)"
if [[ -n "$BASH" ]]; then
  sudo chsh -s "$BASH" "${USER:-$(id -un)}"
fi
sudo apt-get remove -y fonts-powerline powerline zsh || true
rm -rf "$HOME/.local/share/zinit"
