#!/bin/bash

# Claude code
# Check if claude is installed first
if command -v claude >/dev/null; then
  echo "Claude CLI found. Skipping Claude installation"
else
  curl -fsSL https://claude.ai/install.sh | bash
fi

# Copilot CLI
npm install -g @github/copilot

if command -v node >/dev/null; then
  echo "Node found. Skipping node installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "node"
fi

# MCP Hub
npm install -g mcp-hub@latest
