#!/bin/bash

# Taken from https://cloud.google.com/storage/docs/gsutil_install

# curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-441.0.0-linux-x86_64.tar.gz
tar -C "$HOME" -xzf google-cloud-cli-441.0.0-linux-x86_64.tar.gz
~/google-cloud-sdk/install.sh

sudo rm -f google-cloud-sdk-441.0.0-linux-x86_64.tar.gz

