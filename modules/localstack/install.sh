curl --output localstack-cli-4.3.0-linux-amd64-onefile.tar.gz \
  --location https://github.com/localstack/localstack-cli/releases/download/v4.3.0/localstack-cli-4.3.0-linux-amd64-onefile.tar.gz

sudo tar xvzf localstack-cli-4.3.0-linux-*-onefile.tar.gz -C /usr/local/bin

# Remove the downloaded tarball
rm localstack-cli-4.3.0-linux-*-onefile.tar.gz

if [[ ! -d "$HOME/python3/envs/awslocal" ]]; then
  mkdir -p "$HOME"/python3/envs
  cd "$HOME"/python3/envs || exit
  python3 -m venv awslocal
  source "$HOME"/python3/envs/awslocal/bin/activate
  python3 -m pip install awscli-local
  deactivate
fi
