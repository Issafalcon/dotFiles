#!/bin/zsh

# Get the bash completions for azure cli
autoload -U +X compinit && compinit
autoload -U +X bashcompinit && bashcompinit

# TODO: Figure out why this doesn't source properly
source /etc/bash_completion.d/azure-cli
