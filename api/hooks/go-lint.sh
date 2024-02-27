#!/usr/bin/env bash

set -eu -o pipefail

if ! command -v golangci-lint &> /dev/null ; then
    echo "golangci-lint not installed or available in the PATH" >&2
    echo "please check https://github.com/golangci/golangci-lint" >&2
    exit 1
fi

# Check if a go.work file exists in the current directory
if [ -f "go.work" ]; then
    # Iterate through every folder in the current directory
    for dir in */; do
        dir=${dir%/}
        # Check if a go.mod file exists in the folder
        if ls "$dir"/*.go &> /dev/null; then
            echo "Linting module in folder: $(pwd)/$dir"
            # Run lint on that folder
            golangci-lint run ./$dir/... "$@"
        fi
    done
else
    golangci-lint run
    echo "No go.work file found in the current directory"
fi
