#!/bin/bash
sudo apt update
sudo apt install software-properties-common
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt install python3.9 python3-pip python-dev python3-dev python3.8-venv

echo Python version is "$(python3.9 --version)"
