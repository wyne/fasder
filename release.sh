#!/bin/bash

# Check if version is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

# Input version
input_version="$1"

# Get the version from local code
fasder_version=$(go run . -v 2>&1)

# Check if the go_version variable is empty
if [ -z "$fasder_version" ]; then
  echo "'go run . -v' did not return any output. Please ensure the command runs correctly."
  exit 1
fi

# Compare the versions
if [ "$fasder_version" != "$input_version" ]; then
  echo "Version mismatch: fasder version is '$fasder_version', but expected '$input_version'. Exiting."
  exit 1
fi

# Check if the Git tag exists
if ! git rev-parse "$input_version" >/dev/null 2>&1; then
  echo "Git tag '$input_version' does not exist. Please create the tag before proceeding."
  exit 1
fi

# Run the git archive command to create the tar.gz
git archive --prefix=fasder-$input_version/ -o fasder-$input_version.tar.gz main

# Success message
echo "Created fasder-$input_version.tar.gz successfully."
