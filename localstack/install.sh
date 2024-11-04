curl --output localstack-cli-3.8.0-linux-amd64-onefile.tar.gz \
    --location https://github.com/localstack/localstack-cli/releases/download/v3.8.0/localstack-cli-3.8.0-linux-amd64-onefile.tar.gz

sudo tar xvzf localstack-cli-3.8.0-linux-*-onefile.tar.gz -C /usr/local/bin

# Remove the downloaded tarball
rm localstack-cli-3.8.0-linux-*-onefile.tar.gz
