#!/usr/bin/env bash

set -eu -o pipefail

if ! command -v golangci-lint &> /dev/null ; then
    echo "golangci-lint not installed or available in the PATH" >&2
    echo "please check https://github.com/golangci/golangci-lint" >&2
    exit 1
fi

# Check if a api folder exists in the current directory
if [ -d "api" ]; then
    cd api && golangci-lint run ./... "$@"
else
    golangci-lint run ./...
fi
