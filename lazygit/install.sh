#!/bin/bash

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

# Install homebrew
brew --version
if [[ $? -eq 0 ]]; then
  echo "Homebrew found. Skipping Homebrew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

# Lazygit
brew install jesseduffield/lazygit/lazygit
