#!/bin/bash

# Cursor agent CLI
if command -v agent >/dev/null; then
  echo "Cursor CLI found. Skipping Cursor CLI installation"
else
  curl https://cursor.com/install -fsS | bash
fi

rtk init -g --agent cursor
