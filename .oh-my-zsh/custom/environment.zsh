####################
# User configuration
####################

# export MANPATH="/usr/local/man:$MANPATH"

autoload -U +X bashcompinit && bashcompinit

export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

# Generated for envman. Do not edit.
[ -s "$HOME/.config/envman/load.sh" ] && source "$HOME/.config/envman/load.sh"

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

complete -C /home/linuxbrew/.linuxbrew/Cellar/terraform/0.13.4/bin/terraform terraform
complete -o nospace -C /usr/bin/terraform terraform

eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

export EDITOR="nvim"
export ASPNETCORE_ENVIRONMENT=Development
export ASPNETCORE_USE_RUNNING_NODE=true

# set DISPLAY variable to the IP automatically assigned to WSL2
export DISPLAY=$(route.exe print | grep 0.0.0.0 | head -1 | awk '{print $4}'):0.0
