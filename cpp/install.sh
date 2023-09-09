#!/bin/bash
#
# Install clang compiler and required dependencies
sudo apt update && \
  sudo apt install clang \
  gcc \
  g++ \
  make \
  cmake \
  libssl-dev
