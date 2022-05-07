#!/bin/bash

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# Install homebrew
if [[ command -v brew ]]; then
  echo "Homebrew found. Skipping Homebrew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "homebrew"
fi

# Install delta: https://dandavison.github.io/delta/introduction.html
brew install git-delta

# Set delta config
git config --global core.pager "delta --dark --paging=never"
git config --global include.path "~/themes.gitconfig"
git config --global interactive.diffFilter "delta --color-only"
git config --global delta.navigate "true"
git config --global delta.line-numbers "true"
git config --global delta.side-by-side "false"
git config --global delta.syntax-theme "Dracula"
git config --global delta.features "decorations line-numbers zebra-dark"
git config --global merge.conflictstyle "diff3"
 
