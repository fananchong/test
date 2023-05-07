#!/bin/bash

set -ex

TOOLS_DIR=${PWD}

go build

pushd ${PWD}/../go_analysis_test_example
${TOOLS_DIR}/go_analysis_test --path=./...
popd
