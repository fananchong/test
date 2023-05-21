#!/bin/bash

set -ex

TOOLS_DIR=${PWD}

go build

${TOOLS_DIR}/go_analysis_test --path=. --go_module=go_analysis_test
