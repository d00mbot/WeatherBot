#!/usr/bin/env bash
echo "golang-ci lint..."
golangci-lint run ./...

echo "gogroup..."
gogroup -order std,other,prefix=local  $(find . -type f -name "*.go" | grep -v "vendor/")
