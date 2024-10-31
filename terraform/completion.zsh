# Get current version of terragrunt install in homebrew
TG_VERSION=$(terragrunt --version) 
TG_VERSION=${TG_VERSION##* }

autoload -U +X bashcompinit && bashcompinit
complete -o nospace -C /home/linuxbrew/.linuxbrew/Cellar/terragrunt/"$TG_VERSION"/bin/terragrunt terragrunt
