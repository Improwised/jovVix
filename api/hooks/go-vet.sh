#!/bin/bash

set -eu -o pipefail

# Check if a api folder exists in the current directory
if [ -d "api" ]; then
    cd api && go vet ./... "$@"
else
    go vet ./... "$@"
fi
