#!/bin/bash

if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  # set DISPLAY variable to the IP automatically assigned to WSL2
  export DISPLAY=$(route.exe print | grep 0.0.0.0 | head -1 | awk '{print $4}'):0.0
  path+=("/mnt/c/Program Files/Oracle/VirtualBox")
  # Used for vagrant - Enables vagrant use from within WSL2
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
fi
