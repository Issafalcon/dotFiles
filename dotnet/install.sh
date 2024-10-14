#!/bin/bash
# Check if using Ubuntu 22.04
sudo apt update -y &&
  sudo apt install -y dotnet-sdk-8.0
dotnet tool install --global PlantUmlClassDiagramGenerator
