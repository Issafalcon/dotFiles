#!/bin/bash

sudo apt update \
	&& sudo apt install \
		git \
		stow \
		zsh \
		curl \
		wget \
		zip \
		unzip \
		build-essential \
		libssl-dev

# homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

touch "${HOME}/.dotFileModules"
chmod 777 "${HOME}/.dotFileModules"

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

# Need python and pip to install below
"${SCRIPT_DIR}"/node/install.sh

source "$HOME"/.zshrc
