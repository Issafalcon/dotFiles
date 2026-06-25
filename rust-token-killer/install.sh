#!/usr/bin/env bash

# Install brew
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

# https://github.com/rtk-ai/rtk
if command -v rtk >/dev/null; then
  echo "rtk found. Skipping rtk installation"
else
  brew install rtk
fi
