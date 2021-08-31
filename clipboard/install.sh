#!/bin/bash

# Install win yank if on WSL so vim can use windows clipboard
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
	echo "Windows 10 Bash"
	echo "Fetching win32yank to use windows clipboard in vim"

	curl -sLo/tmp/win32yank.zip https://github.com/equalsraf/win32yank/releases/download/v0.0.4/win32yank-x64.zip
	unzip -p /tmp/win32yank.zip win32yank.exe >/tmp/win32yank.exe
	chmod +x /tmp/win32yank.exe
	
	mv /tmp/win32yank.exe "$HOME"/.local/bin
else
	sudo apt-get install -y xclip
fi
