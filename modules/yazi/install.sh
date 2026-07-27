#!/bin/bash
# Idempotent yazi install via Homebrew (+ optional projects.yazi plugin).
set -euo pipefail

SCRIPT_DIR=$(cd "${0%/*}" && pwd -P)

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
"$brew_bin" install yazi ffmpeg fd

plugin_dir="${SCRIPT_DIR}/.config/yazi/plugins/projects.yazi"
if [[ -d "$plugin_dir/.git" ]]; then
  echo "projects.yazi already present; skipping clone"
elif [[ -e "$plugin_dir" ]]; then
  echo "projects.yazi path exists but is not a git repo; leaving as-is"
else
  mkdir -p "$(dirname "$plugin_dir")"
  git clone https://github.com/MasouShizuka/projects.yazi.git "$plugin_dir"
fi
