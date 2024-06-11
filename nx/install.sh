#!/bin/bash

# Update the apt package index and install packages needed to use the nx globally
sudo apt update
sudo apt install -y jq

npm add --global nx@latest
