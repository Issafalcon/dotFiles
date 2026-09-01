#!/bin/bash
# https://docs.astral.sh/uv/getting-started/installation#uninstallation
set -euo pipefail

export PATH="${HOME}/.local/bin:${HOME}/.cargo/bin:${PATH}"

if command -v uv >/dev/null 2>&1; then
  uv cache clean || true
  # ponytail: dirs may be absent on a fresh/minimal install
  py_dir="$(uv python dir 2>/dev/null || true)"
  tool_dir="$(uv tool dir 2>/dev/null || true)"
  [[ -n "${py_dir}" && -d "${py_dir}" ]] && rm -rf "${py_dir}"
  [[ -n "${tool_dir}" && -d "${tool_dir}" ]] && rm -rf "${tool_dir}"
fi

rm -f "${HOME}/.local/bin/uv" "${HOME}/.local/bin/uvx"
# Pre-0.5.0 installs lived under ~/.cargo/bin
rm -f "${HOME}/.cargo/bin/uv" "${HOME}/.cargo/bin/uvx"

echo "uv uninstalled"
