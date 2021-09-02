#!/bin/bash

# Taken from https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-using-native-package-management

# Update the apt package index and install packages needed to use the Kubernetes apt repository:
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl

# Add kubernetes apt repo
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# Download the Google Cloud public signing key:
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

# Update apt package index with new repo and install kubectl
sudo apt-get update
sudo apt-get install -y kubectl

# Install k9s tool
curl -sS https://webinstall.dev/k9s | bash

