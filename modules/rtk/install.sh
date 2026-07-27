#!/usr/bin/env bash
SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  echo "homebrew module is required. Install it from the TUI first." >&2
  exit 1
fi

# https://github.com/rtk-ai/rtk
if command -v rtk >/dev/null; then
  echo "rtk found. Skipping rtk installation"
else
  brew install rtk
fi
