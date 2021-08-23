#!/bin/bash
DOTFILES_ROOT=$(pwd -P)

ARGS=$(getopt -a --options m:iur --long "module:,install,uninstall,remove" -- "$@")
MODULE=""
INSTALL=false
UNINSTALL=false
REMOVE=false

eval set -- "$ARGS"

set -e

while true; do
  case "$1" in
    -m|--module)
 			MODULE="${2}"
      shift 2;;
    -i|--install)
      INSTALL=true
      shift;;
    -u|--uninstall)
      UNINSTALL=true
      shift;;
    -r|--remove)
      REMOVE=true
      shift;;
    --)
      break;;
  esac
done

echo ''

info() {
	# shellcheck disable=SC2059
	printf "\r  [ \033[00;34m..\033[0m ] $1\n"
}

user() {
	# shellcheck disable=SC2059
	printf "\r  [ \033[0;33m??\033[0m ] $1\n"
}

success() {
	# shellcheck disable=SC2059
	printf "\r\033[2K  [ \033[00;32mOK\033[0m ] $1\n"
}

fail() {
	# shellcheck disable=SC2059
	printf "\r\033[2K  [\033[0;31mFAIL\033[0m] $1\n"
	echo ''
	exit
}

find_zsh() {
	if command -v zsh >/dev/null 2>&1 && grep "$(command -v zsh)" /etc/shells >/dev/null; then
		command -v zsh
	else
		echo "/bin/zsh"
	fi
}

update_module_config() {
	if [[ $REMOVE ]]; then
		sed -i "/$MODULE/d" ~/.dotFileModules
	else
		if ! grep -q "${MODULE}" "$HOME"/.dotFileModules; then
			echo "${MODULE}" >> "${HOME}"/.dotFileModules
		fi
	fi
}

stow_module() {
	stow "${MODULE}"
}

install_module_dependencies() {
	./"${MODULE}"/install.sh
	if [[ $? -eq 0 ]]; then
		success "Successfully installed dependencies for $MODULE module"
	else
		fail "Failed to install dependencies for $MODULE module"
	fi
	check last status code
}

if [[ "${MODULE}" -eq "zsh" ]]; then
	ZSH="$(find_zsh)"
	test "$(expr "$SHELL" : '.*/\(.*\)')" != "ZSH" &&
	command -v chsh >/dev/null 2>&1 &&
	chsh -s "$ZSH" &&
	success "set $("$ZSH" --version) at $ZSH as default shell"
fi

if [[ $INSTALL ]]; then
	install_module_dependencies
fi

add_module_config
stow_module
