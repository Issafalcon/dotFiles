#!/bin/bash

# Install libsecret to store git credentials
if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # Get windows current user home dir
  WINDOWS_HOME=$(wslpath "$(cmd.exe /C "echo %USERPROFILE%" 2>/dev/null | tr -d '\r')")

  # Dbus UI not available on WSL. Use wincred store instead
  git config --global credential.helper "$WINDOWS_HOME/AppData/Local/Programs/Git/mingw64/bin/git-credential-manager.exe"
else
  sudo apt install gnome-keyring
  sudo apt-get install libsecret-1-0 libsecret-1-dev
  cd /usr/share/doc/git/contrib/credential/libsecret
  sudo make
  git config --global credential.helper /usr/share/doc/git/contrib/credential/libsecret/git-credential-libsecret
fi
