#!/bin/bash
set -euo pipefail

sudo rm -f /usr/local/bin/localstack
rm -rf "$HOME/python3/envs/awslocal"
