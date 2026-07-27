#!/bin/bash
# Idempotent rustup install. Prefer -y so streaming (non-TTY) installs work;
# with requires_input: true the TUI can also run this interactively via ExecProcess.
set -euo pipefail

if command -v rustup >/dev/null 2>&1; then
  echo "rustup already installed; updating stable"
else
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
fi

# shellcheck disable=SC1091
if [[ -f "${HOME}/.cargo/env" ]]; then
  # ponytail: source for this process; interactive shells use ~/.cargo/env from profile
  source "${HOME}/.cargo/env"
fi

rustup update stable
rustc --version
