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

  echo "Installing plugins..."

  copilot plugin marketplace add DietrichGebert/ponytail
  copilot plugin install ponytail@ponytail
fi

# Install rtk
if command -v rtk >/dev/null; then
  echo "rtk found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "rtk"
fi

rtk init -g
