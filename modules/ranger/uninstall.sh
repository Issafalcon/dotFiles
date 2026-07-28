#!/bin/bash
set -euo pipefail

sudo apt-get remove -y ranger xsel
rm -rf ~/.config/ranger/plugins/ranger_devicons
