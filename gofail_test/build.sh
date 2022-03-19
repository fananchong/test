#!/bin/bash

set -ex

export PATH=$(go env GOPATH)/bin:$PATH
n=$(which gofail | wc -l)
if [ "$n" == "0" ]; then
    go get github.com/etcd-io/gofail
    go mod tidy
fi

pushd examples && gofail enable && popd
go mod tidy
go build
