#!/bin/bash

source "$DOTFILES/azure/completion.zsh"

TENANTS=$(az account tenant list | jq -r 'map(.tenantId) | join(" ")')
# convert TENANTS json array to bash array
TENANTS_LIST=($TENANTS)

echo "Select a tenant to login to: "

select TENANT in "${TENANTS_LIST[@]}"; do
	az login -t "$TENANT" --allow-no-subscriptions
	break;
done
