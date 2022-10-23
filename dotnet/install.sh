#!/bin/bash
# Check if using Ubuntu 22.04
if [ "$(lsb_release -rs)" = "22.04" ]; then
	sudo apt update &&
		sudo apt install dotnet-sdk-6.0
  else
    echo "Not on Ubuntu 22.04. Please refer to https://learn.microsoft.com/en-us/dotnet/core/install/linux-ubuntu?source=recommendations"
fi
