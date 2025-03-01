#!/usr/bin/env bash

if [ -f .env ]; then
  source .env
else
  echo "Warning: .env file not found."
fi

if [ -z "$GITHUB_TOKEN" ]; then
  echo "Error: GITHUB_TOKEN variable is not set."
  exit 1
fi

GITHUB_REPO_PATH=squaredLabs

set_github_repo_token() {
  git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/${GITHUB_REPO_PATH}".insteadOf "https://github.com/${GITHUB_REPO_PATH}"
}

unset_github_repo_token() {
  git config --global --unset url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/${GITHUB_REPO_PATH}".insteadOf
}

# unset the token if the first argument is "unset"
if [ "$1" = "unset" ]; then
  unset_github_repo_token
  echo "Token auth unset."
  exit 0
fi


# by default, set the token
set_github_repo_token
echo "Token auth set."
