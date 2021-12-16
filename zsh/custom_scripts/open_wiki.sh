#!/bin/bash

ARGS=$(getopt -a --options w:h --long "wikiName:,help" -- "$@")
WIKI=""

usage()
{
	printf "Usage guide: \n"
	printf " Arguments:\n"
	printf "	--wikiName	The name of the wiki guide to open"
}

eval set -- "$ARGS"

set -e

while true; do
	case "$1" in
		-w | --wikiName)
			WIKI="${2}"
			shift 2
			;;
		-h| --help)
			usage
			break
			;;
		--)
			break
			;;
	esac
done

openWiki() {
	WIKIDIR="${PROJECTS}/wiki"

	if [[ -n $WIKI ]]; then
		find "${WIKIDIR}" -maxdepth 5 -type f -iname "*${WIKI}*.pdf" -exec zathura '{}' \;
	fi
}

openWiki
