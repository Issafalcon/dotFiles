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
		-m | --module)
			MODULE="${2}"
			shift 2
			;;
		-i | --install)
			INSTALL=true
			shift
			;;
		-u | --uninstall)
			UNINSTALL=true
			shift
			;;
		-r | --remove)
			REMOVE=true
			shift
			;;
		--)
			break
			;;
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

find_shell() {
	if command -v "${1}" >/dev/null 2>&1 && grep "$(command -v "${1}")" /etc/shells >/dev/null; then
		command -v "$1"
	else
		echo "/bin/$1"
	fi
}

update_module_config() {
	if [[ $REMOVE == true ]]; then

		sed -i "/$MODULE/d" ~/.dotFileModules \
			&& stow -D "${MODULE}"

		if [[ $? -eq 0 ]]; then
			success "Successfully removed and unstowed $MODULE module"
		else
			fail "Failed to remove and unstow $MODULE module"
		fi
		
		# Swap shell back to bash if we are removing zsh module
		if [[ ${MODULE} == "zsh" ]]; then
			BASH="$(find_shell bash)"
			command -v chsh >/dev/null 2>&1 \
				&& chsh -s "$BASH" \
				&& success "set $("$BASH" --version) at $BASH as default shell"
		fi
	else
		if (! grep -q "${MODULE}" "$HOME"/.dotFileModules) && [[ "${MODULE}" != "zsh" ]]; then
			echo "${MODULE}" >>"${HOME}"/.dotFileModules
		fi

		stow "${MODULE}"

		if [[ $? -eq 0 ]]; then
			success "Successfully stowed $MODULE module"
		else
			fail "Failed to stow $MODULE module"
		fi
	fi
}

install_module_dependencies() {
	if [[ -f ./"${MODULE}"/install.sh ]]; then
		./"${MODULE}"/install.sh
		if [[ $? -eq 0 ]]; then

			success "Successfully installed dependencies for $MODULE module"
		else
			fail "Failed to install dependencies for $MODULE module"
		fi
	fi
}

if [[ ${MODULE} == "zsh" ]]; then
	ZSH="$(find_shell zsh)"
	test "$(expr "$SHELL" : '.*/\(.*\)')" != "ZSH" \
		&& command -v chsh >/dev/null 2>&1 \
		&& chsh -s "$ZSH" \
		&& success "set $("$ZSH" --version) at $ZSH as default shell"
fi

if [[ $INSTALL == true ]]; then
	install_module_dependencies
fi

update_module_config

source "$HOME"/.zshrc
