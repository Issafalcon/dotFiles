#!/bin/bash

sudo apt-get update &&
	sudo apt-get install \
		texlive-full \
		latexmk

sudo apt update
sudo apt install xindy

pip install pygments
