#!/bin/bash
# Transfers all pdf files in the wiki repo to a target directory, maintaining folder structure
# Args:
#     $1: Taget directory root

if [[ $# -lt 1 ]]; then
  echo "$0: Target directory is required"
  exit 2
fi

# Transfers all .pdf files in a directory recursively
find "$PROJECTS"/wiki -name '*.pdf' -exec cp '{}' "$1" \;
