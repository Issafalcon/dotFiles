#!/bin/bash

# Need python and pip to install below
if command -v python3 >/dev/null; then
  echo "Python 3 found. Skipping python 3 installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "python"
  path+=(/usr/bin/pip3)
fi

if [[ ! -d "$HOME/python3/envs/qmk" ]]; then
  mkdir -p "$HOME"/python3/envs
  cd "$HOME"/python3/envs || exit
  python3 -m venv qmk
  source "$HOME"/python3/envs/qmk/bin/activate
  python3 -m pip install qmk
  deactivate
fi

cd "$HOME"/python3/envs/qmk || exit
source ./bin/activate
qmk setup -H "$PROJECTS/qmk_firmware" Issafalcon/qmk_firmware
qmk config user.keyboard=keebart/sofle_choc_pro
qmk config user.keymap=issafalcon
deactivate
