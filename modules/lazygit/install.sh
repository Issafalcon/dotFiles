#!/bin/bash

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

# Install homebrew
brew --version
if [[ $? -eq 0 ]]; then
  echo "Homebrew found. Skipping Homebrew installation"
else
	echo "homebrew module is required. Install it from the TUI first." >&2
	exit 1
fi

# Lazygit
brew install jesseduffield/lazygit/lazygit
