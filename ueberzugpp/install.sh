#!/bin/bash

# Install brew if not present
if command -v brew >/dev/null; then
  echo "brew found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

brew install jstkdng/programs/ueberzugpp
