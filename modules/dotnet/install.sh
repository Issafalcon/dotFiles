#!/bin/bash
sudo apt update -y

sudo add-apt-repository ppa:dotnet/backports

# Ask to install .net pass in version number as argument
install_dotnet() {
  if [ -z "$1" ]; then
    echo "Please provide the .net version to install as an argument."
    exit 1
  fi
  # Ask if the user wants to install .net
  read -rp "Do you want to install .net $1? (y/n) " answer
  if [ "$answer" != "y" ]; then
    echo "Skipping .net $1 installation."
    return
  fi

  sudo apt install -y dotnet-sdk-"$1"
}

install_dotnet 8.0
install_dotnet 9.0
install_dotnet 10.0

dotnet tool install --global PlantUmlClassDiagramGenerator
dotnet tool install --global lazydotnet
dotnet tool install --global dotnet-outdated-tool
