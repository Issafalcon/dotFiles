#!/bin/bash
sudo apt update
sudo apt install python3
sudo apt install python3-pip 
sudo apt install python3-dev
sudo apt install python3-venv

echo Python version is "$(python3 --version)"
