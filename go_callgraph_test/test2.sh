#!/bin/bash

set -ex

go build

TOOLDIR=${PWD}

pushd /data/fananchong/torchlight/git/backend/src/haidao/backend/game
${TOOLDIR}/go_callgraph_test -algo static .
popd
