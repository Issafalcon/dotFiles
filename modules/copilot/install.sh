#!/bin/bash
SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

if command -v node >/dev/null; then
  echo "Node found. Skipping node installation"
else
  echo "node module is required. Install it from the TUI first." >&2
  exit 1
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
  echo "rtk module is required. Install it from the TUI first." >&2
  exit 1
fi

rtk init -g
