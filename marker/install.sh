#!/bin/bash
# Marker install script
# Installs marker-pdf into a Python virtualenv for direct host execution.
# This avoids Docker memory overhead which causes OOM on ML workloads.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
MARKER_VENV="${XDG_DATA_HOME:-$HOME/.local/share}/marker/venv"
OLLAMA_DEFAULT_MODEL="llama3.2"

# +----------+
# | SWAP     |
# +----------+

# Marker's ML models need ~12-14GB RSS on CPU (float32 is hardcoded for CPU).
# A swapfile ensures the kernel can handle peak usage without OOM-killing the process.
_marker_ensure_swap() {
  local required_gb=16
  local swapfile="/swapfile-marker"

  local total_swap_kb
  total_swap_kb=$(awk '/SwapTotal/ {print $2}' /proc/meminfo)
  local total_swap_gb=$(( total_swap_kb / 1024 / 1024 ))

  if (( total_swap_gb >= required_gb )); then
    echo "✓ Swap is sufficient: ${total_swap_gb}GB available."
    return 0
  fi

  echo "Current swap: ${total_swap_gb}GB. Adding ${swapfile} (${required_gb}GB) for marker model loading..."
  sudo fallocate -l "${required_gb}G" "${swapfile}"
  sudo chmod 600 "${swapfile}"
  sudo mkswap "${swapfile}"
  sudo swapon "${swapfile}"

  # Persist across reboots
  if ! grep -q "${swapfile}" /etc/fstab; then
    echo "${swapfile} none swap sw 0 0" | sudo tee -a /etc/fstab
  fi

  echo "✓ Swapfile created: $(free -h | awk '/Swap/ {print $2}') total swap now available."
}

_marker_ensure_swap

# +----------+
# | PYTHON   |
# +----------+

# Prefer python3.11; marker supports 3.10+
PYTHON_BIN=""
for candidate in python3.11 python3.12 python3.10 python3; do
  if command -v "${candidate}" >/dev/null 2>&1; then
    version=$("${candidate}" -c 'import sys; print(sys.version_info[:2])')
    if "${candidate}" -c 'import sys; sys.exit(0 if sys.version_info >= (3,10) else 1)' 2>/dev/null; then
      PYTHON_BIN="${candidate}"
      echo "Using ${PYTHON_BIN}: $(${PYTHON_BIN} --version)"
      break
    fi
  fi
done

if [[ -z "${PYTHON_BIN}" ]]; then
  echo "Python 3.10+ is required. Installing python3.11..."
  sudo apt-get update && sudo apt-get install -y python3.11 python3.11-venv python3.11-dev
  PYTHON_BIN="python3.11"
fi

# +-----------+
# | VENV      |
# +-----------+

if [[ -f "${MARKER_VENV}/bin/marker_single" ]]; then
  echo "Marker already installed at ${MARKER_VENV}"
  echo "To upgrade: ${MARKER_VENV}/bin/pip install -U marker-pdf[full]"
else
  echo ""
  echo "Creating virtualenv at ${MARKER_VENV}..."
  mkdir -p "$(dirname "${MARKER_VENV}")"
  "${PYTHON_BIN}" -m venv "${MARKER_VENV}"

  echo "Installing PyTorch (CPU)..."
  "${MARKER_VENV}/bin/pip" install --quiet --upgrade pip
  "${MARKER_VENV}/bin/pip" install --quiet torch torchvision \
    --index-url https://download.pytorch.org/whl/cpu

  echo "Installing marker-pdf..."
  "${MARKER_VENV}/bin/pip" install --quiet "marker-pdf[full]"

  echo "✓ Marker installed at ${MARKER_VENV}"
fi

# +---------+
# | OLLAMA  |
# +---------+

echo ""
echo "Ollama provides free local LLMs to enhance marker's accuracy (--use-llm mode)."
read -r -p "Install Ollama for local LLM support? [Y/n] " install_ollama
install_ollama="${install_ollama:-Y}"

if [[ "${install_ollama,,}" == "y" ]]; then
  if command -v ollama >/dev/null 2>&1; then
    echo "Ollama already installed: $(ollama --version)"
  else
    echo "Installing Ollama..."
    curl -fsSL https://ollama.com/install.sh | bash
  fi

  # Install systemd resource limits to prevent Ollama freezing the system
  if systemctl list-units --type=service 2>/dev/null | grep -q ollama; then
    OLLAMA_DROPIN_DIR="/etc/systemd/system/ollama.service.d"
    echo "Installing Ollama resource limits..."
    sudo mkdir -p "${OLLAMA_DROPIN_DIR}"
    sudo cp "${SCRIPT_DIR}/ollama-limits.conf" "${OLLAMA_DROPIN_DIR}/limits.conf"
    sudo systemctl daemon-reload
    sudo systemctl restart ollama 2>/dev/null || true
    echo "✓ Ollama memory limits applied. Edit ${OLLAMA_DROPIN_DIR}/limits.conf to adjust."
  fi

  echo ""
  read -r -p "Pull default model '${OLLAMA_DEFAULT_MODEL}'? [Y/n] " pull_model
  pull_model="${pull_model:-Y}"
  if [[ "${pull_model,,}" == "y" ]]; then
    ollama pull "${OLLAMA_DEFAULT_MODEL}"
    echo "✓ Model '${OLLAMA_DEFAULT_MODEL}' ready."
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
echo "Run 'marker_convert --help' for all options."
