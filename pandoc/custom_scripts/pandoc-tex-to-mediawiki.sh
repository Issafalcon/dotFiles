#!/bin/bash

# Script to convert a LaTeX wiki file to mediawiki format using pandoc
# Requires pandoc and pandoc-citeproc to be installed
# If a bibliography.bib file is present, it will also use that to create citations using the elsevier-with-titles.csl style
# Usage: pandoc-tex-to-mediawiki.sh <file.tex>

if [[ $# -lt 1 ]]; then
  echo 1>$2 "$0: File argument is required"
  exit 2
fi

FULL_PATH=$(realpath "$1")

FILE="$(basename "${FULL_PATH}")"
TEX_DIR="$(dirname "${FULL_PATH}")"

if ! command -v pandoc &>/dev/null; then
  echo "pandoc could not be found"
  exit 1
fi

if ! command -v pandoc-citeproc &>/dev/null; then
  echo "pandoc-citeproc could not be found"
  exit 1
fi

# Check if mediawiki folder exists, if not, create it
if ! [ -d "${TEX_DIR}/../mediawiki/" ]; then
  mkdir "${TEX_DIR}/../mediawiki"
fi

cd "$TEX_DIR" || return

if [ -f "${TEX_DIR}/../bibliography.bib" ]; then
  pandoc -f latex -t mediawiki \
    --metadata link-citations=true \
    --bibliography="${TEX_DIR}/../bibliography.bib" \
    --csl="${PROJECTS}/wiki/elsevier-with-titles.csl" \
    "${FULL_PATH}" \
    -o "${TEX_DIR}/../mediawiki/${FILE}.md"
else
  pandoc -f latex -t mediawiki \
    --metadata link-citations=true \
    --csl="${PROJECTS}/wiki/elsevier-with-titles.csl" \
    "${FULL_PATH}" \
    -o "${TEX_DIR}/../mediawiki/${FILE}.md"
fi
