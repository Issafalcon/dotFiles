#!/bin/bash
# Idempotent ImageMagick install via Homebrew.
set -euo pipefail

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
  echo "brew not found; install (or re-run) the homebrew module first" >&2
  exit 1
fi

eval "$("$brew_bin" shellenv)"
"$brew_bin" install imagemagick
