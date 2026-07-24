#!/bin/bash

# homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

sudo apt-get install build-essential
/home/linuxbrew/.linuxbrew/bin/brew install gcc
