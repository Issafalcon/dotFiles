#!/bin/bash

sudo apt update \
  && sudo add-apt-repository ppa:git-core/ppa \
	&& sudo apt install \
		git \
		stow \
		zsh \
		curl \
		wget \
		zip \
		unzip \
		build-essential \
		libssl-dev \
    jq

touch "${HOME}/.dotFileModules"
chmod 777 "${HOME}/.dotFileModules"

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

"${SCRIPT_DIR}"/bootstrap.sh "-i" "-m" "node"

/bin/bash "${SCRIPT_DIR}"/bootstrap.sh "-i" "-m" "zsh"
