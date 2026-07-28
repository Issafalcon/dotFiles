if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  path+=(/mnt/c/Program Files/Oracle/VirtualBox)
fi
