#!/bin/bash

sudo rm -rf /usr/local/jira_1.1.0_linux_x86_64
wget https://github.com/ankitpokhrel/jira-cli/releases/download/v1.1.0/jira_1.1.0_linux_x86_64.tar.gz
sudo tar -C /usr/local/ -xzf jira_1.1.0_linux_x86_64.tar.gz
rm -f jira_1.1.0_linux_x86_64.tar.gz
