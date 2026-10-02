#!/bin/bash

sudo apt-get update -y &&
  sudo apt-get install -y openjdk-17-jdk

echo "Java version is $(java --version)"
