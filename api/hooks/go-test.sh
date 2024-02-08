#!/bin/bash

rm cover.out
gotestsum --format pkgname -- -coverprofile=cover.out ./api/...
