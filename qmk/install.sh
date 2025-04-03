#!/bin/bash

# Need python and pip to install below
if command -v python3 >/dev/null; then
  echo "Python 3 found. Skipping python 3 installation"
else
  "${SCRIPT_DIR}"/../bootstrap.sh "-i" "-m" "python"
  path+=(/usr/bin/pip3)
fi

echo "$1"

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

# Use command line argument to detemine whether to use vial
if [[ "$1" == "vial" ]]; then
  qmk setup -H "$PROJECTS/vial-qmk" -b vial vial-kb/vial-qmk
  cd "$PROJECTS/vial-qmk" || exit
  git submodule update --init --recursive
else
  qmk setup -H "$PROJECTS/qmk_firmware" qmk/qmk_firmware
fi

qmk config user.keyboard=keebart/sofle_choc_pro
qmk config user.keymap=issafalcon

git clone https://github.com/Issafalcon/qmk_userspace.git "${PROJECTS}/qmk_userspace"
qmk config user.overlay_dir="$(realpath "${PROJECTS}/qmk_userspace")"

deactivate
