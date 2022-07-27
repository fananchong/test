#!/bin/bash

set -ex

go build .
docker build -t go_coredump_test -f ./Dockerfile .
docker rm -f test1 || true
docker run --name=test1 \
    --ulimit core=-1 --security-opt seccomp=unconfined \
    -v ${PWD}/run:/myworkdir \
    go_coredump_test
