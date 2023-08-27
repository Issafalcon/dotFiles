#!/bin/bash

curl -LO https://github.com/wez/wezterm/releases/download/20230712-072601-f4abf8fd/wezterm-20230712-072601-f4abf8fd.Ubuntu22.04.deb
sudo apt install -y ./wezterm-20230712-072601-f4abf8fd.Ubuntu22.04.deb

rm wezterm-20230712-072601-f4abf8fd.Ubuntu22.04.deb

