#!/bin/bash

SCRIPT_DIR=$( cd ${0%/*} && pwd -P )

# Install latest aws cli
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "${SCRIPT_DIR}/awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install


rm -f "${SCRIPT_DIR}"/awscliv2.zip
rm -fR "${SCRIPT_DIR}"/dist
rm -f "${SCRIPT_DIR}"/install
rm -f "${SCRIPT_DIR}"/README.md
rm -f "${SCRIPT_DIR}"/THIRD_PARTY_LICENSES
