#!/bin/bash
SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Install brew if not present
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  echo "homebrew module is required. Install it from the TUI first." >&2
  exit 1
fi

brew install jstkdng/programs/ueberzugpp
