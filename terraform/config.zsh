# Install the oh-my-zsh terraform plugin as a zinit snippet
zinit snippet OMZ::plugins/terraform/terraform.plugin.zsh

# Add completions from the plugin
zinit ice as"completion"
zinit snippet OMZ::plugins/terraform/_terraform
