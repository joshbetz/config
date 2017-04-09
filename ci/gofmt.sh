#!/usr/bin/env bash

fmt="$(find . -type f -name '*.go' -print0 | xargs -0 gofmt -l)"

if [ -n "$fmt" ]; then
	echo "Unformatted source code:"
	echo "$fmt"
	exit 1
fi
