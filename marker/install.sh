#!/bin/bash
# Marker install script
# Installs the marker document-to-markdown tool as a Docker image,
# and optionally installs Ollama for free local LLM-enhanced conversions.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
MARKER_IMAGE_NAME="marker-pdf"
OLLAMA_DEFAULT_MODEL="llama3.2"

# +---------+
# | DOCKER  |
# +---------+

# Ensure Docker Engine (dockerd) is installed — not just the CLI
if ! command -v dockerd >/dev/null 2>&1; then
  echo "Docker Engine (dockerd) not found. Installing docker-ce..."
  if ! apt-cache show docker-ce >/dev/null 2>&1; then
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
      | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
      https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
      | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
    sudo apt-get update
  fi
  sudo apt-get install -y docker-ce
  sudo usermod -aG docker "$USER"
  echo "NOTE: Log out and back in (or run 'newgrp docker') for Docker group membership to take effect."
else
  echo "Docker found: $(docker --version)"
fi

# Ensure Docker daemon is running
_marker_ensure_docker() {
  if docker info >/dev/null 2>&1; then
    return 0
  fi

  echo "Docker daemon is not running. Attempting to start it..."
  if sudo systemctl start docker 2>/dev/null; then
    : # systemctl succeeded
  elif sudo service docker start 2>/dev/null; then
    : # service succeeded
  else
    echo "Starting dockerd directly..."
    sudo dockerd &>/tmp/dockerd.log &
  fi

  # Wait up to 20 seconds for the daemon to become ready
  local attempts=0
  while ! docker info >/dev/null 2>&1; do
    if (( attempts++ >= 20 )); then
      echo "ERROR: Docker daemon did not start in time." >&2
      echo "Try running: sudo dockerd &" >&2
      return 1
    fi
    sleep 1
  done
  echo "Docker daemon started."
}

_marker_ensure_docker || exit 1

echo ""
echo "Building marker Docker image (this may take several minutes on first run)..."
docker build -t "${MARKER_IMAGE_NAME}" "${SCRIPT_DIR}"

echo ""
echo "✓ Marker Docker image '${MARKER_IMAGE_NAME}' built successfully."

# +------------------+
# | MODELS CACHE DIR |
# +------------------+

MARKER_MODELS_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/marker/models"
mkdir -p "${MARKER_MODELS_DIR}"
echo "✓ Model cache directory: ${MARKER_MODELS_DIR}"

# +---------+
# | OLLAMA  |
# +---------+

echo ""
echo "Ollama provides free local LLMs to enhance marker's accuracy (--use_llm mode)."
read -r -p "Install Ollama for local LLM support? [Y/n] " install_ollama
install_ollama="${install_ollama:-Y}"

if [[ "${install_ollama,,}" == "y" ]]; then
  if command -v ollama >/dev/null 2>&1; then
    echo "Ollama already installed: $(ollama --version)"
  else
    echo "Installing Ollama..."
    curl -fsSL https://ollama.com/install.sh | bash
  fi

  echo ""
  read -r -p "Pull default model '${OLLAMA_DEFAULT_MODEL}'? [Y/n] " pull_model
  pull_model="${pull_model:-Y}"
  if [[ "${pull_model,,}" == "y" ]]; then
    ollama pull "${OLLAMA_DEFAULT_MODEL}"
    echo "✓ Model '${OLLAMA_DEFAULT_MODEL}' ready."
  fi

  # Start Ollama service if not running
  if ! ollama list >/dev/null 2>&1; then
    echo "Starting Ollama service..."
    ollama serve &>/dev/null &
    sleep 2
  fi
fi

# +---------+
# | SUMMARY |
# +---------+

echo ""
echo "================================================"
echo " Marker installation complete!"
echo "================================================"
echo ""
echo "Usage (via the marker_convert shell function):"
echo ""
echo "  # Basic conversion (PDF/EPUB → markdown):"
echo "  marker_convert /path/to/book.pdf"
echo ""
echo "  # With local LLM enhancement (requires Ollama):"
echo "  marker_convert /path/to/book.pdf --use-llm"
echo ""
echo "  # Specify output directory:"
echo "  marker_convert /path/to/book.epub ~/Notes/book-notes"
echo ""
echo "  # Force OCR (for scanned PDFs):"
echo "  marker_convert /path/to/scanned.pdf --force-ocr"
echo ""
echo "  # Using GitHub Copilot as LLM backend:"
echo "  marker_convert /path/to/book.pdf --use-llm --llm-service copilot"
echo ""
echo "Run 'marker_convert --help' for all options."
