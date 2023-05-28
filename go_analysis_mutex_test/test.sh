#!/bin/bash

set -ex

TOOLS_DIR=${PWD}

go build

${TOOLS_DIR}/go_analysis_mutex_test --path=${PWD}/test5
