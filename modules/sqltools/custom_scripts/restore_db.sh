#!/bin/bash

usage() {
  # Print out the usage instructions for this script
  echo "Usage: $0 [options]"
  echo ""
  echo "Options:"
  echo "  -n, --name            Container name"
  echo "  -p, --password        SQL Server SA password"
  echo "  -b, --backup-path     Path to backup file"
  echo "  -d, --database-name   Name of database to restore into"
  echo ""
  echo "Examples:"
  echo "  $0 -n sqlserver -p 'MyPassword' -b /home/user/backup.bak -d my-database"
  echo ""
  exit 1
}

if [ -z "$1" ]; then
	echo "No argument supplied"
	exit 1
fi

ARGS=$(getopt -a --options n:p:b:d:h --long "name:password:backup-path:database-name:help" -- "$@")
CONTAINER_NAME=""
SQL_SA_PASSWORD=""
BACKUP_FILE_PATH=""
DB_NAME=""

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
	-b | --backup-path)
		BACKUP_FILE_PATH="${2}"
		shift 2
		;;
	-d | --database-name)
		DB_NAME="${2}"
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

if [ -z "$BACKUP_FILE_PATH" ]; then
  echo "Missing required argument: -b | --backup-path"
  usage
fi

if [ -z "$DB_NAME" ]; then
  echo "Missing required argument: -d | --database-name"
  usage
fi

# Get just the filename and extension from BACKUP_FILE_PATH
BACKUP_FILE_NAME=$(basename "$BACKUP_FILE_PATH")
RESTORE_QUERY="RESTORE DATABASE $DB_NAME FROM DISK = '/var/opt/mssql/backup/$BACKUP_FILE_NAME'"

echo "Copying backup file to container..."
docker exec -it "$CONTAINER_NAME" mkdir -p /var/opt/mssql/backup
docker cp "$BACKUP_FILE_PATH" "$CONTAINER_NAME":/var/opt/mssql/backup

echo "Restoring database..."
BACKUP_FILE_OUTPUT=$(docker exec -it "$CONTAINER_NAME" /opt/mssql-tools/bin/sqlcmd -S localhost \
	-U SA -P "$SQL_SA_PASSWORD" \
	-Q "RESTORE FILELISTONLY FROM DISK = '/var/opt/mssql/backup/$BACKUP_FILE_NAME'")

MDF_OUTPUT=$(echo "$BACKUP_FILE_OUTPUT" | { grep -i -E ".*\.mdf" || true; })

if [ -z "$MDF_OUTPUT" ]; then
	echo "No MDF file found in backup"
	exit 1
else
	echo "MDF file found in backup"
	MDF_FILE_NAME=$(echo "$MDF_OUTPUT" | tr -s ' ' | awk '{print $1;}')
	MDF_FILE_PATH=$(echo "$MDF_OUTPUT" | rev | cut -d'\' -f 1 | rev | awk '{print $1;}')
	RESTORE_QUERY="$RESTORE_QUERY WITH MOVE '$MDF_FILE_NAME' TO '/var/opt/mssql/data/$MDF_FILE_PATH'"
	echo "$RESTORE_QUERY"
fi

NDF_OUTPUT=$(echo "$BACKUP_FILE_OUTPUT" | { grep -i -E ".*\.ndf" || true; })

if [ -z "$NDF_OUTPUT" ]; then
	echo "No NDF file found in backup"
else
	echo "NDF file found in backup"
	NDF_FILE_NAME=$(echo "$NDF_OUTPUT" | tr -s ' ' | awk '{print $1;}')
	NDF_FILE_PATH=$(echo "$NDF_OUTPUT" | rev | cut -d'\' -f 1 | rev | awk '{print $1;}')
	RESTORE_QUERY="$RESTORE_QUERY, MOVE '$NDF_FILE_NAME' TO '/var/opt/mssql/data/$NDF_FILE_PATH'"
fi

LDF_OUTPUT=$(echo "$BACKUP_FILE_OUTPUT" | { grep -i -E ".*\.ldf" || true; })

if [ -z "$LDF_OUTPUT" ]; then
	echo "No LDF file found in backup"
else
	echo "LDF file found in backup"
	LDF_FILE_NAME=$(echo "$LDF_OUTPUT" | tr -s ' ' | awk '{print $1;}')
	LDF_FILE_PATH=$(echo "$LDF_OUTPUT" | rev | cut -d'\' -f 1 | rev | awk '{print $1;}')
	RESTORE_QUERY="$RESTORE_QUERY, MOVE '$LDF_FILE_NAME' TO '/var/opt/mssql/data/$LDF_FILE_PATH'"
fi

echo "Executing restore query..."
echo "$RESTORE_QUERY"

docker exec -it "$CONTAINER_NAME" /opt/mssql-tools/bin/sqlcmd \
	-S localhost -U SA -P "$SQL_SA_PASSWORD" \
	-Q "$RESTORE_QUERY"
