#!/bin/bash

# Check if cursor is installed first
if command -v agent >/dev/null; then
  echo "Cursor CLI found. Skipping Cursor CLI installation"
else
  curl https://cursor.com/install -fsS | bash
fi

# Install rtk
if command -v rtk >/dev/null; then
  echo "rtk found. Skipping brew installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "rtk"
fi

rtk init -g --agent cursor
