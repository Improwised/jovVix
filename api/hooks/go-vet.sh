#!/bin/bash

set -eu -o pipefail

# Check if a go.work file exists in the current directory
if [ -f "go.work" ]; then
    # Iterate through every folder in the current directory
    for dir in */; do
        dir=${dir%/}
        # Check if a go.mod file exists in the folder
        if ls "$dir"/*.go &> /dev/null; then
            echo "Go module in folder: $(pwd)/$dir"
            # Run lint on that folder
            go vet ./$dir/... "$@"
        fi
    done
else
    go vet ./... "$@"
    echo "No go.work file found in the current directory"
fi
