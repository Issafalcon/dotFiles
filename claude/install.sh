#!/bin/bash

# Claude code
# Check if claude is installed first
if command -v claude >/dev/null; then
  echo "Claude CLI found. Skipping Claude installation"
else
  curl -fsSL https://claude.ai/install.sh | bash
fi
