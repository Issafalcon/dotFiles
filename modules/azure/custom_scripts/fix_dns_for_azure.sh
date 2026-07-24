#!/bin/bash

command sudo bash -c "sed -i '/management.azure.com/d' /etc/hosts" ; sudo bash -c 'echo "$(dig management.azure.com | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}$") management.azure.com" >> /etc/hosts'
