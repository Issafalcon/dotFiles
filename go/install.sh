#!/bin/bash

sudo rm -rf /usr/local/go
wget https://dl.google.com/go/go1.17.6.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.17.6.linux-amd64.tar.gz
rm -f go1.17.6.linux-amd64.tar.gz
 
