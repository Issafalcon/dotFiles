#!/bin/bash

# Install libsecret to store git credentials
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
	# Dbus UI not available on WSL. Use wincred store instead
	git config --global credential.helper "/mnt/c/Program\ Files/Git/mingw64/bin/git-credential-manager-core.exe"
else
	sudo apt install gnome-keyring
	sudo apt-get install libsecret-1-0 libsecret-1-dev
	cd /usr/share/doc/git/contrib/credential/libsecret
	sudo make
	git config --global credential.helper /usr/share/doc/git/contrib/credential/libsecret/git-credential-libsecret
fi
