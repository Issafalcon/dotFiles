#!/bin/bash

usage() {
  # Print out the usage instructions for this script
  echo "Usage: $0 [options]"
  echo ""
  echo "Options:"
  echo "  -n, --name        Container name"
  echo "  -p, --password    SQL Server SA password"
  echo "  -v, --volume      Local volume to mount. Accepts an absolute file path to a directory on local machine (requires read/write permissions for docker))"
  echo "                    or a docker volume name (e.g. my-vol). If the volume does not exist, it will be created."
  echo "  --port            (Optional) Port to expose on host. Default: 1433"
  echo ""
  echo "Examples:"
  echo "  $0 -n sqlserver -p 'MyPassword' -v /home/user/sqlserver"
  echo "  $0 -n sqlserver -p 'MyPassword' -v my-docker-volume"
  echo ""
  exit 1
}

ARGS=$(getopt -a --options n:p:v:h --long "name:password:volume:help:,port:" -- "$@")
CONTAINER_NAME=""
SQL_SA_PASSWORD=""
LOCAL_VOLUME=""
PORT="1433"

eval set -- "$ARGS"

set -e

while true; do
  case "$1" in
  -n | --name)
    CONTAINER_NAME="${2}"
    shift 2
    ;;
  -p | --password)
    SQL_SA_PASSWORD="${2}"
    shift 2
    ;;
  -v | --volume)
    LOCAL_VOLUME="${2}"
    shift 2
    ;;
  --port)
    PORT="${2}"
    shift 2
    ;;
  -h | --help)
    usage
    ;;
  --)
    break
    ;;
  esac
done

if [ -z "$CONTAINER_NAME" ]; then
  echo "Missing required argument: -n | --name"
  usage
fi

if [ -z "$SQL_SA_PASSWORD" ]; then
  echo "Missing required argument: -p | --password"
  usage
fi

if [ -z "$LOCAL_VOLUME" ]; then
  echo "Missing required argument: -v | --volume"
  usage
fi

docker run -e 'ACCEPT_EULA=Y' -e "MSSQL_SA_PASSWORD=$SQL_SA_PASSWORD" \
  --name "$CONTAINER_NAME" -p "$PORT":1433 \
  -v "$LOCAL_VOLUME":/var/opt/mssql \
  -d mcr.microsoft.com/mssql/server:2022-latest
