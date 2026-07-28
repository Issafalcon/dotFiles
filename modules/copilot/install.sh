#!/bin/bash

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

rtk init -g
