#!/bin/bash
# Standalone installer: https://docs.astral.sh/uv/getting-started/installation
# UV_NO_MODIFY_PATH: path comes from this module's path.zsh, not shell profile edits.
set -euo pipefail

export UV_NO_MODIFY_PATH=1

if command -v uv >/dev/null 2>&1; then
  echo "uv already installed; self-updating"
  uv self update
else
  curl -LsSf https://astral.sh/uv/install.sh | sh
fi

# Ensure this process can find uv before a new shell loads path.zsh
export PATH="${HOME}/.local/bin:${PATH}"

uv --version
