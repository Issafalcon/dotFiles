#!/bin/bash

sudo apt-get update \
 && sudo apt-get install tmux 

# Tmux Plugin manager and plugins
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm

