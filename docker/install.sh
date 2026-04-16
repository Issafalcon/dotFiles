#!/bin/bash
# Docker Engine install script for Ubuntu
# https://docs.docker.com/engine/install/ubuntu/

set -euo pipefail

# Skip if Docker Engine (dockerd) is already fully installed
if command -v dockerd >/dev/null 2>&1; then
  echo "Docker Engine already installed: $(docker --version)"
  exit 0
fi

# Remove unofficial/conflicting packages
echo "Removing any conflicting Docker packages..."
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
  sudo apt-get remove -y "$pkg" 2>/dev/null || true
done

# Add Docker's official GPG key
sudo apt-get update
sudo apt-get install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the Docker apt repository
# Use UBUNTU_CODENAME for derivatives (e.g. Linux Mint); fall back to VERSION_CODENAME
DOCKER_CODENAME=$(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  ${DOCKER_CODENAME} stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null

# Install Docker Engine, CLI, containerd, and plugins
sudo apt-get update
sudo apt-get install -y \
  docker-ce \
  docker-ce-cli \
  containerd.io \
  docker-buildx-plugin \
  docker-compose-plugin

# Post-install: allow running docker without sudo
sudo groupadd -f docker
sudo usermod -aG docker "$USER"

# Enable and start Docker on boot
sudo systemctl enable docker.service
sudo systemctl enable containerd.service
sudo systemctl start docker.service

echo ""
echo "✓ Docker Engine installed: $(docker --version)"
echo "NOTE: Log out and back in (or run 'newgrp docker') to use Docker without sudo."
