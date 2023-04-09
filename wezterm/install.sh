#!/bin/bash

curl -LO https://github.com/wez/wezterm/releases/download/20230326-111934-3666303c/wezterm-20230326-111934-3666303c.Ubuntu22.04.deb
sudo apt install -y ./wezterm-20230326-111934-3666303c.Ubuntu22.04.deb

if [ ! -d ~/bin ]; then
    mkdir ~/bin
fi

mv ./WezTerm-20230408-112425-69ae8472-Ubuntu20.04.AppImage ~/bin/wezterm
~/bin/wezterm
