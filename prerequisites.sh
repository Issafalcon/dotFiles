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

# Install nvm and node
"${SCRIPT_DIR}"/node/install.sh

/bin/bash "${SCRIPT_DIR}"/bootstrap.sh "-i" "-m" "zsh"
