#!/bin/bash

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

# Install nvm and node
"${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"

# Lazygit
brew install jesseduffield/lazygit/lazygit
