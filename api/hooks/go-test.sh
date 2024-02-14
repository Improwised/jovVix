#!/bin/bash

set -eu -o pipefail

if ! command -v gotestsum &> /dev/null ; then
    echo "gotestsum not installed or available in the PATH" >&2
    echo "please check https://github.com/gotestyourself/gotestsum" >&2
    exit 1
fi

# Check if a go.work file exists in the current directory
if [ -f "go.work" ]; then
    # Iterate through every folder in the current directory
    for dir in */; do
        dir=${dir%/}
        # Check if a go.mod file exists in the folder
        if ls "$dir"/*.go &> /dev/null; then
            echo "Testing module in folder: $(pwd)/$dir"
            # Run lint on that folder
            source ./$dir/.env.testing && gotestsum --format pkgname -- -coverprofile=cover.out ./$dir/... "$@"
        fi
    done
else
    echo "No go.work file found in the current directory"
fi
