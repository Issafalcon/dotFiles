if grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null; then
  export VAGRANT_WSL_ENABLE_WINDOWS_ACCESS="1"
fi
