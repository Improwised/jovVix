#!/bin/bash

set -eu -o pipefail

if ! command -v gotestsum &> /dev/null ; then
    echo "gotestsum not installed or available in the PATH" >&2
    echo "please check https://github.com/gotestyourself/gotestsum" >&2
    exit 1
fi

# Check if a api folder exists in the current directory
if [ -d "api" ]; then
    cd api && source .env.testing && gotestsum --format pkgname -- -coverprofile=cover.out ./... "$@"
else
    source .env.testing && gotestsum --format pkgname -- -coverprofile=cover.out ./... "$@"
fi
