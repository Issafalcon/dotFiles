#!/bin/bash
set -euo pipefail

# Keep sudo warm before the long brew installer (credential expires ~15m).
sudo apt-get update -y
sudo apt-get install -y build-essential

# Non-interactive brew install (no "Press RETURN" prompt under the TUI).
export NONINTERACTIVE=1
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

BREW="/home/linuxbrew/.linuxbrew/bin/brew"
if [[ ! -x "$BREW" ]]; then
  BREW="/opt/homebrew/bin/brew"
fi
if [[ ! -x "$BREW" ]]; then
  BREW="$(command -v brew)"
fi
"$BREW" install gcc
