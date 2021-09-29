#!/bin/bash

sudo apt-get update && sudo apt-get install ranger xsel

# Ranger plugins for dev-icons
git clone https://github.com/alexanderjeurissen/ranger_devicons ~/.config/ranger/plugins/ranger_devicons
