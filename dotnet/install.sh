#!/bin/bash
# Install individual components due to issue with package mapping
# see https://github.com/dotnet/sdk/issues/27082#issuecomment-1211143446 for details
sudo apt install \
  aspnetcore-runtime-6.0=6.0.8-1 \
  dotnet-apphost-pack-6.0=6.0.8-1 \
  dotnet-host=6.0.8-1 \
  dotnet-hostfxr-6.0=6.0.8-1 \
  dotnet-runtime-6.0=6.0.8-1 \
  dotnet-sdk-6.0=6.0.400-1 \
  dotnet-targeting-pack-6.0=6.0.8-1
