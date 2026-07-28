#!/bin/bash

curl -LO https://github.com/quarto-dev/quarto-cli/releases/download/v1.8.27/quarto-1.8.27-linux-amd64.deb
sudo apt install -y ./quarto-1.8.27-linux-amd64.deb
rm quarto-1.8.27-linux-amd64.deb
