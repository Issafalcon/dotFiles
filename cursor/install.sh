#!/bin/bash

# Claude code
# Check if claude is installed first
if command -v agent >/dev/null; then
  echo "Cursor CLI found. Skipping Cursor CLI installation"
else
  curl https://cursor.com/install -fsS | bash
fi
