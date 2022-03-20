#!/bin/bash

set -ex

export PATH=$(go env GOPATH)/bin:$PATH
n=$(which gofail | wc -l)
if [ "$n" == "0" ]; then
    go get github.com/etcd-io/gofail
    go mod tidy
fi

gofail enable examples
go mod tidy
go build
