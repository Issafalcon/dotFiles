#!/bin/bash

curl -LO https://github.com/wez/wezterm/releases/download/nightly/wezterm-nightly.Ubuntu22.04.deb
sudo apt install -y ./wezterm-nightly.Ubuntu22.04.deb

rm wezterm-nightly.Ubuntu22.04.deb
