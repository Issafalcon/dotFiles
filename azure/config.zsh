#!/bin/zsh

# Get the bash completions for azure cli
source /etc/bash_completion.d/azure-cli

# Install the oh-my-zsh az cli plugin as a zinit snippet
zinit load dmakeienko/azcli
