#!/bin/bash

# Claude code
# Check if claude is installed first
if command -v claude >/dev/null; then
  echo "Claude CLI found. Skipping Claude installation"
else
  curl -fsSL https://claude.ai/install.sh | bash
fi

# Register MCP servers at user scope (global across all projects)
if command -v claude >/dev/null; then
  echo "Registering MCP servers..."

  claude mcp add --scope user nvim-mcp -- nvim-mcp --log-file . --log-level debug --connect auto
  claude mcp add --scope user context7 -- npx -y @upstash/context7-mcp@latest

  echo "MCP servers registered."
fi

# Add additinal useful skills
if command -v claude >/dev/null; then
  echo "Adding additional skills..."
  claude plugin install superpowers@claude-plugins-official

  # Add the pptx-posters skill from the scientific-agent-skills repository (will install the skills manager client if not already installed)
  npx skills add https://github.com/k-dense-ai/scientific-agent-skills --skill pptx-posters

  echo "Additional skills added."
fi
