#!/usr/bin/env bash

# The following can have several formats
# v1.1.3-1-g1397b0b-dirty
# aka v + latest-tag + '-' + #commits-above-tag + '-g' + short-commit-id + '-dirty'
# the final '-dirty' suffix can be empty
#
# v1.1.3 when on master tagged branch

# Remove leading v
GIT_REF=$(git describe --always --tags --dirty | sed -e 's/^v//')

if [ -z "$1" ] || [ "$1" == "version" ]; then
	# shellcheck disable=SC1001
	VERSION=$(echo "$GIT_REF" | awk -F\- '{print $1}')
	if [ -n "$VERSION" ]; then
		echo "$VERSION"
		exit 0
	fi
fi

if [ "$1" == "revision" ]; then
	# shellcheck disable=SC1001
	VERSION=$(echo "$GIT_REF" | awk -F\- 'NF>3{printf "%s-%s-%s\n",$2,$3,$4}NF==3{printf "%s-%s\n",$2,$3}')
	if [ -n "$VERSION" ]; then
		echo "$VERSION"
		exit 0
	fi
	git rev-parse --short HEAD
	exit 0
fi

echo "n/a"
