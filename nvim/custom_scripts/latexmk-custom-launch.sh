#!/bin/bash

# Script to manually execute latexmk compiler that also compiles indexes and glossaries
# This is to workaround the limitation that VimTex \ll command is unable to also compile the glossaries
# and will error is one is included

if [[ $# -lt 1 ]]; then
  echo 1>$2 "$0: File argument is required"
  exit 2
fi

FULL_PATH=$(realpath "$1")

FILE="$(basename "${FULL_PATH}")"
TEX_DIR="$(dirname "${FULL_PATH}")"

if ! command -v latexmk &> /dev/null; then
 echo "latexmk could not be found" 
 exit 1
fi

cd "$TEX_DIR" || return
latexmk -pdf -pdflatex=lualatex -verbose -shell-escape "${FILE}"

if ! command -v makeglossaries &> /dev/null; then
 echo "makeglossaries could not be found" 
 exit 1
fi

makeglossaries "${FILE%.*}"
latexmk -pdf -pdflatex=lualatex -pvc -verbose -file-line-error -synctex=1 -interaction=nonstopmode -shell-escape "${FILE}"
