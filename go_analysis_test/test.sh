#!/bin/bash

set -ex

TOOLS_DIR=${PWD}

go build

${TOOLS_DIR}/go_analysis_test --path=${PWD}/../go_analysis_test_example/
