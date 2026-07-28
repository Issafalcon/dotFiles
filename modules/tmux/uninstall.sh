#!/bin/bash
set -euo pipefail

sudo apt-get remove -y tmux
rm -rf ~/.tmux/plugins/tpm
