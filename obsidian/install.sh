#!/bin/bash

sudo apt update -y
sudo apt install libnss3-dev
sudo curl -Lo /usr/bin/obsidian https://github.com/obsidianmd/obsidian-releases/releases/download/v1.8.9/Obsidian-1.8.9.AppImage
sudo chmod 777 /usr/bin/obsidian
