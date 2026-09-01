#!/bin/bash
# https://github.com/Graphify-Labs/graphify — PyPI package is graphifyy; CLI is graphify
set -euo pipefail

export PATH="${HOME}/.local/bin:${PATH}"

if ! command -v uv >/dev/null 2>&1; then
  echo "uv not found; install the 'uv' module first" >&2
  exit 1
fi

uv tool install graphifyy

# Cursor writes .cursor/rules/ relative to cwd; run from $HOME for user-global rules.
# Claude/Copilot skill dirs are already under ~/.claude and ~/.copilot.
cd "${HOME}"

register() {
  local platform="$1"
  echo "Registering graphify for ${platform}..."
  graphify install --platform "${platform}"
}

found=0

if command -v claude >/dev/null 2>&1 || [[ -d "${HOME}/.claude" ]]; then
  register claude
  found=1
else
  echo "Claude Code not detected — skipping"
fi

if command -v agent >/dev/null 2>&1 || [[ -d "${HOME}/.cursor" ]]; then
  register cursor
  found=1
else
  echo "Cursor not detected — skipping"
fi

if command -v copilot >/dev/null 2>&1 || [[ -d "${HOME}/.copilot" ]]; then
  register copilot
  found=1
else
  echo "GitHub Copilot CLI not detected — skipping"
fi

if [[ "${found}" -eq 0 ]]; then
  echo "No AI tools detected (claude / cursor / copilot)."
  echo "CLI is installed; later: graphify install --platform <name>"
fi

graphify --version
