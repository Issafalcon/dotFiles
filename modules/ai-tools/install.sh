#!/bin/bash

# Claude code
# Check if claude is installed first
if command -v claude >/dev/null; then
  echo "Claude CLI found. Skipping Claude installation"
else
  curl -fsSL https://claude.ai/install.sh | bash
fi

# Copilot CLI
if command -v copilot >/dev/null; then
  echo "Copilot CLI found. Skipping Copilot installation"
else
  echo "Installing Copilot CLI..."
  npm install -g @github/copilot
fi

if command -v node >/dev/null; then
  echo "Node found. Skipping node installation"
else
	echo "node module is required. Install it from the TUI first." >&2
	exit 1
fi

# MCP Hub
npm install -g mcp-hub@latest
