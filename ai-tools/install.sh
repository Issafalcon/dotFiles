#!/bin/bash

# Claude code
curl -fsSL https://claude.ai/install.sh | bash

# Copilot CLI
npm install -g @github/copilot

if command -v node >/dev/null; then
  echo "Node found. Skipping node installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "node"
fi

# MCP Hub
npm install -g mcp-hub@latest
