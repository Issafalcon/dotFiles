#!/bin/bash

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

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
  echo "rtk module is required. Install it from the TUI first." >&2
  exit 1
fi

rtk init -g --agent cursor
