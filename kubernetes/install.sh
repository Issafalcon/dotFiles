#!/bin/bash

# Taken from https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-using-native-package-management

# Update the apt package index and install packages needed to use the Kubernetes apt repository:
sudo apt-get update
# apt-transport-https may be a dummy package; if so, you can skip that package
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg

# If the folder `/etc/apt/keyrings` does not exist, it should be created before the curl command, read the note below.
# sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg # allow unprivileged APT programs to read this keyring
#
# This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list # helps tools such as command-not-found to work correctly

# Update apt package index with new repo and install kubectl
sudo apt-get update
sudo apt-get install -y kubectl

# Install k9s tool
curl -sS https://webinstall.dev/k9s | bash

# Load kubectl completions here into local_function as sourcing the kubectl doesn't seem to work
# source <(kubectl completion zsh) - Doesn't work
# Set the kubectl completion code for zsh[1] to autoload on startup
kubectl completion zsh >"$HOME/zsh_local/functions/_kubectl"
