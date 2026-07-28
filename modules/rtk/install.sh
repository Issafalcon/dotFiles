#!/usr/bin/env bash
SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

# https://github.com/rtk-ai/rtk
if command -v rtk >/dev/null; then
  echo "rtk found. Skipping rtk installation"
else
  brew install rtk
fi
