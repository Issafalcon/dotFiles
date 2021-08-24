#!/bin/bash

# Install nvm and node
rm -rf "$NVM_DIR" ~/.npm ~/.bower # Remove existing nvm and node
cd "$HOME"/

git clone https://github.com/nvm-sh/nvm.git .nvm
cd "$HOME"/.nvm
git checkout v0.38.0

export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm

nvm install --lts
nvm use --lts
