#!/bin/bash

function create_release {
  local tag_name="v$1"
  local http_response=$(curl \
    -u $GITHUB_USER_TOKEN \
    -w "HTTP_STATUS:%{http_code}" -X POST -s \
    -H "Content-Type: application/json" \
    -d "{\"tag_name\":\"$tag_name\", \"target_commitish\":\"master\", \"name\":\"cli release $tag_name\", \"prerelease\": false, \"draft\": false}" \
    "https://api.github.com/repos/squarescale/squarescale-cli/releases")

  local http_body=$(echo $http_response | sed -e 's/HTTP_STATUS\:.*//g')
  local http_status=$(echo $http_response | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

  if [ ! $http_status -eq 201 ]; then
    echo "Error [HTTP status: $http_status]" >&2
    echo $http_body >&2
    exit 1
  fi

  echo $http_body | sed -re 's/.*"url":\s*".*\/([0-9]*)".*/\1/'
}

function push_release_data {
  local id=$1
  local name=$2
  local path="$(pwd)/$name"

  local http_response=$(curl                    \
    -u $GITHUB_USER_TOKEN                       \
    -w "HTTP_STATUS:%{http_code}" -X POST       \
    -s                                          \
    -H "Content-Type: application/octet-stream" \
    --data-binary "@$path"                      \
    "https://uploads.github.com/repos/squarescale/squarescale-cli/releases/$id/assets?name=$name")

  local http_body=$(echo $http_response | sed -e 's/HTTP_STATUS\:.*//g')
  local http_status=$(echo $http_response | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

  if [ ! $http_status -eq 201 ]; then
    echo "Error [HTTP status: $http_status]" >&2
    echo $http_body >&2
    exit 1
  fi
}

set -e

echo "Create release on Github..."
version=$(grep -r "const Version string = \"*\"" * | awk '{ print $5 }' | sed 's/"//g')
release_id=$(create_release $version)

for exe in "sqsc-linux-amd64" "sqsc-darwin-amd64"
do
  echo "Push $exe to Github..."
  push_release_data $release_id $exe
done

echo "Done."
