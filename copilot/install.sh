#!/bin/bash

if command -v node >/dev/null; then
  echo "Node found. Skipping node installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "node"
fi

# Copilot CLI
if command -v copilot >/dev/null; then
  echo "Copilot CLI found. Skipping Copilot installation"
else
  echo "Installing Copilot CLI..."
  npm install -g @github/copilot
fi
