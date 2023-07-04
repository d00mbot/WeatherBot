#!/usr/bin/env bash

function get_dependencies() {
    declare -a packages=(
        "github.com/Bubblyworld/gogroup/..."
    )

    ## now loop through the above array
    for pkg in "${packages[@]}"
    do
        echo "$pkg"
        go get -u -v "$pkg"
    done

    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    golangci-lint --version

}

echo Gonna to update go tools and packages...
get_dependencies
echo All is done!