# Install the oh-my-zsh docker and docker compose plugins as a zinit snippet

# Just the completion snippets for docker (the OMZ plugin tries to load this snippet file from the ZSH_CACHE location, which won't be found)
zinit snippet https://github.com/docker/cli/blob/master/contrib/completion/zsh/_docker

# Docker compose plugin (with all the aliases) and the completion script for docker-compose
zinit snippet OMZ::plugins/docker-compose/docker-compose.plugin.zsh
zinit snippet OMZ::plugins/docker-compose/_docker-compose
