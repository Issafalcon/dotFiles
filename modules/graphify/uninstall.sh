#!/bin/bash
# https://github.com/Graphify-Labs/graphify
set -euo pipefail

export PATH="${HOME}/.local/bin:${PATH}"

# Cursor rule path is cwd-relative; uninstall from $HOME so ~/.cursor/rules is cleared.
cd "${HOME}"

if command -v graphify >/dev/null 2>&1; then
  for platform in claude cursor copilot; do
    graphify "${platform}" uninstall 2>/dev/null || true
  done
fi

if command -v uv >/dev/null 2>&1; then
  uv tool uninstall graphifyy 2>/dev/null || true
else
  rm -f "${HOME}/.local/bin/graphify" "${HOME}/.local/bin/graphify-mcp"
fi

echo "graphify uninstalled"
