#!/bin/bash
# Check minimum version of ubuntu 20.04
sudo apt update -y &&
  sudo apt install -y dotnet-sdk-8.0

dotnet tool install --global PlantUmlClassDiagramGenerator

if [ $(lsb_release -rs) == "20.04" ]; then
  sudo add-apt-repository ppa:dotnet/backports
fi

sudo apt install -y dotnet-sdk-9.0

dotnet tool install --global lazydotnet
dotnet tool install --global dotnet-outdated-tool
