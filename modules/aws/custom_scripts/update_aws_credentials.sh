#!/bin/bash

export AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION:-"eu-west-1"}

# We need to have at least one cache file with SSO credentials.
SSO_CACHE_FILES=$(ls -t "$HOME"/.aws/sso/cache/*.json 2>/dev/null || true)
if [ -n "${SSO_CACHE_FILES}" ] && [ -z "${AWS_SESSION_TOKEN}" ]; then
	echo "Exporting AWS SSO/role credentials to environment variables"

	SSO_ROLE=$(aws sts get-caller-identity --query=Arn | cut -d '_' -f 2)
	SSO_ACCOUNT=$(aws sts get-caller-identity --query=Account --output text)
	SESSION_FILE=$(echo $SSO_CACHE_FILES | tr ' ' '\n' | head -n 1)
	SSO_ACCESS_TOKEN=$(jq -r '.accessToken' "$SESSION_FILE")
	SSO_REGION=$(jq -r '.region' "$SESSION_FILE")
	AWS_CREDENTIALS=$(aws sso get-role-credentials \
		--role-name="$SSO_ROLE" --account-id="$SSO_ACCOUNT" \
		--access-token="$SSO_ACCESS_TOKEN" --region="$SSO_REGION")

	LOCAL_AWS_ACCESS_KEY_ID=$(echo "$AWS_CREDENTIALS" | jq -r '.roleCredentials.accessKeyId')
	LOCAL_AWS_SECRET_ACCESS_KEY=$(echo "$AWS_CREDENTIALS" | jq -r '.roleCredentials.secretAccessKey')
	LOCAL_AWS_SESSION_TOKEN=$(echo "$AWS_CREDENTIALS" | jq -r '.roleCredentials.sessionToken')

	if [ -s "${HOME}/.aws/credentials" ]; then
		sed -i -e "s|aws_access_key_id=.*|aws_access_key_id=${LOCAL_AWS_ACCESS_KEY_ID}|" \
			-e "s|aws_secret_access_key=.*|aws_secret_access_key=${LOCAL_AWS_SECRET_ACCESS_KEY}|" \
			-e "s|aws_session_token=.*|aws_session_token=${LOCAL_AWS_SESSION_TOKEN}|" \
			"${HOME}"/.aws/credentials
	fi

	export AWS_ACCESS_KEY_ID=$LOCAL_AWS_ACCESS_KEY_ID
	export AWS_SECRET_ACCESS_KEY=$LOCAL_AWS_SECRET_ACCESS_KEY
	export AWS_SESSION_TOKEN=$LOCAL_AWS_SESSION_TOKEN
fi
