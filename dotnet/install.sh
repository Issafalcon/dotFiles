#!/bin/bash

# Download and install the automated dotnet install script from mcr to install .NET 5.0
wget https://packages.microsoft.com/config/ubuntu/21.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
sudo dpkg -i packages-microsoft-prod.deb
rm packages-microsoft-prod.deb

sudo apt-get update
sudo apt-get install -y apt-transport-https \
	&& sudo apt-get update \
	&& sudo apt-get install -y dotnet-sdk-6.0

dotnet tool install --global xunit-cli --version 0.1.16

