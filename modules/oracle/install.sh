#!/bin/bash

SCRIPT_DIR=$(cd ${0%/*} && pwd -P)

python3 -m venv ~/oci-cli-env
source ~/oci-cli-env/bin/activate
pip install oci-cli
