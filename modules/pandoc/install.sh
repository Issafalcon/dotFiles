#!/bin/bash

sudo apt update -y &&
  sudo apt install -y \
    pandoc

# For pandoc to pdf conversion, we need to install the following packages
# This will install pdflatex and required fonts
sudo apt-get install texlive-latex-base
sudo apt-get install texlive-fonts-recommended
sudo apt-get install texlive-fonts-extra
sudo apt-get install texlive-latex-extra
