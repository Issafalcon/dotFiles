#!/bin/bash
set -euo pipefail
SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

brew_bin="$(command -v brew 2>/dev/null || true)"
if [[ -z "$brew_bin" ]]; then
  for candidate in /home/linuxbrew/.linuxbrew/bin/brew /opt/homebrew/bin/brew /usr/local/bin/brew; do
    if [[ -x "$candidate" ]]; then
      brew_bin="$candidate"
      break
    fi
  done
fi
if [[ -z "$brew_bin" ]]; then
  echo "brew not found; install the homebrew module first" >&2
  exit 1
fi

"$brew_bin" install imagemagick
