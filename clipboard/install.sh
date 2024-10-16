#!/bin/bash

# Install win yank if on WSL so vim can use windows clipboard
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  echo "Windows 10 Bash"
  echo "Fetching win32yank to use windows clipboard in vim"

  curl -sLo/tmp/win32yank.zip https://github.com/equalsraf/win32yank/releases/download/v0.0.4/win32yank-x64.zip
  unzip -p /tmp/win32yank.zip win32yank.exe >/tmp/win32yank.exe
  chmod +x /tmp/win32yank.exe

   if [[ ! -d "$HOME/.local/bin" ]]; then
     mkdir "$HOME"/.local/bin
   fi

  sudo mv /tmp/win32yank.exe "$HOME"/.local/bin
  echo "Windows 11 Bash"
  echo "Setting path of win32yank to point to default scoop install location"
  echo "Not installing on WSL"
else
  sudo apt-get install -y xclip
fi
