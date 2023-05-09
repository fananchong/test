#!/bin/bash

set -ex

go build

TOOLDIR=${PWD}

pushd ../go_analysis_test_example/app1/
${TOOLDIR}/go_callgraph_test -algo vta .
popd
