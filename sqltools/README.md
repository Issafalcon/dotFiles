# SQL Server Docker

## Scripts

There are two scripts to help run SQL Server in Docker containers:
1. `run_db.sh`
2. `restore_db.sh`

### Pre-requisites

Bash or another POSIX compliant shell needs to be installed on the machine to run the scripts.
- Anywhere that `git` is installed on Windows, bash will also be installed as a required dependency
- bash is the default shell on many Linux based distros and comes pre-installed

Docker needs to be installed on the machine, and be executable by the user who is running the scripts.
- For Windows, [Docker Desktop](https://www.docker.com/products/docker-desktop/) can be installed for a more user friendly interface
- Once installed Docker Desktop will need to be configured to use Linux containers (this is the default for new installs)

### `run_db.sh`

This script will run a SQL Server 2022 instance in a Linux container
- If the image doesn't already exist locally, it will pull the image automatically
- The default port to expose to the host will be 1433 (default SQL server port)
- To spin up multiple instances of SQL Server, each container will need a unique name, and will need to expose a unique port number on the host
  - This will enable connections from the host to the container via different ports (i.e. Via SSMS - See below)

### `restore_db.sh`

This script will restore a SQL backup into a **running** SQL Server container, so it will be necessary to start a container before running this script
- The backup file will need to be a locally stored backup, and not be restricted by permissions
- You will need the same container name and SQL_SA_PASSWORD as arguments, that were used in the initial `run_db.sh` script to create the container
- Once the backup is restored, you will be able to access the SQL data files in the location of the mounted volume, provided to the container setup in `run_db.sh`

// TODO:
- Additional instructions on connecting to docker container
- Troubleshooting guide
