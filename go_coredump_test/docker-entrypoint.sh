#!/bin/sh

set -ex

cp ${TMPDIR}/* ${WORKDIR}/
${WORKDIR}/go_coredump_test
