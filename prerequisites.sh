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
