#!/bin/bash

INSTALL=false
UNINSTALL=false
REMOVE=false
ALL=false

while getopts ":iura" option; do
	case $option in
		i)
			INSTALL=true
			;;
		u)
			UNINSTALL=true
			;;
		r)
			REMOVE=true
			;;
		a)
			ALL=true
			;;
		\?) # Invalid option
			echo "Error: Invalid option"
			exit
			;;
	esac
done

shift $((OPTIND - 1))

run_module_bootstrap() {
	if [[ $INSTALL == true ]]; then
		./bootstrap.sh -i -m "$1"
	elif [[ $UNINSTALL == true ]]; then
		./bootstrap.sh -u -m "$1"
	elif [[ $REMOVE == true ]]; then
		./bootstrap.sh -r -m "$1"
	else
		./bootstrap.sh -m "$1"
	fi
}

if [[ $ALL == true ]]; then
	for dir in ./*; do
		if [[ $dir != ".git" && $dir != "docs" ]]; then
			run_module_bootstrap "$dir"
		fi
	done
else
	for module in "$@"; do
		run_module_bootstrap "$module"
	done
fi
