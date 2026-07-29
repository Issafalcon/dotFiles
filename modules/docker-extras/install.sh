#!/bin/bash

set -euo pipefail

# Post-install: allow running docker without sudo
sudo groupadd -f docker
sudo usermod -aG docker "$USER"
