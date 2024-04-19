#!/bin/bash

sudo rm -rf /usr/local/go
wget https://dl.pstmn.io/download/latest/linux_64/ -O postman-linux-x64.tar.gz
sudo tar -C /usr/local/ -xzf postman-linux-x64.tar.gz
rm -f postman-linux-x64.tar.gz
 
