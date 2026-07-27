#!/bin/bash
# Homebrew on Linux / WSL2 — idempotent (safe to re-run).
# Requirements: https://docs.brew.sh/Homebrew-on-Linux
set -euo pipefail

install_system_devtools() {
  # System C compiler + standard development tools (brew's own gcc does NOT
  # replace /usr/bin/cc for bootstrap and formula post-install steps).
  if command -v apt-get >/dev/null 2>&1; then
    sudo apt-get update -y
    sudo apt-get install -y build-essential procps curl file git
  elif command -v dnf >/dev/null 2>&1; then
    sudo dnf group install -y development-tools || sudo dnf group install -y "Development Tools"
    sudo dnf install -y procps-ng curl file git
  elif command -v pacman >/dev/null 2>&1; then
    sudo pacman -Sy --noconfirm --needed base-devel procps-ng curl file git
  else
    echo "Unsupported package manager; install a system C toolchain manually." >&2
    exit 1
  fi
}

resolve_brew() {
  for candidate in \
    /home/linuxbrew/.linuxbrew/bin/brew \
    /opt/homebrew/bin/brew \
    /usr/local/bin/brew
  do
    if [[ -x "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
  done
  if command -v brew >/dev/null 2>&1; then
    command -v brew
    return 0
  fi
  return 1
}

ensure_shellenv() {
  local brew_bin="$1"
  local marker="# Homebrew (dotfiles-tui)"
  local line
  line="eval \"\$(${brew_bin} shellenv)\""

  for rc in "${HOME}/.bashrc" "${HOME}/.zshrc"; do
    touch "$rc"
    if grep -Fq 'linuxbrew/.linuxbrew' "$rc" 2>/dev/null \
      || grep -Fq 'brew shellenv' "$rc" 2>/dev/null; then
      continue
    fi
    {
      echo ""
      echo "$marker"
      echo "$line"
    } >>"$rc"
    echo "Appended brew shellenv to $rc"
  done
}

install_system_devtools

BREW="$(resolve_brew || true)"
if [[ -z "${BREW}" ]]; then
  echo "Installing Homebrew…"
  export NONINTERACTIVE=1
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  BREW="$(resolve_brew)"
fi

if [[ -z "${BREW}" || ! -x "${BREW}" ]]; then
  echo "brew not found after install" >&2
  exit 1
fi

echo "Using brew at ${BREW}"
eval "$("${BREW}" shellenv)"
ensure_shellenv "${BREW}"

# Homebrew-provided gcc is still useful for some formulae; safe if already installed.
"${BREW}" install gcc

echo "Homebrew ready: $("${BREW}" --version | head -n1)"
