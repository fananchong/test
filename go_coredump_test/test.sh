#!/bin/bash

set -ex

rm -rf core.*

ulimit -c unlimited

# go build -gcflags="all=-N -l" .

go build .

GOTRACEBACK=crash ./go_coredump_test

# pid=$(ps -ef | grep go_coredump_test | grep -v grep | awk -F' ' '{print $2}')
# Ctrl+\ or kill -SIGQUIT ${pid} or gcore ${pid}

# dlv core ./go_coredump_test core.xxxx
# goroutines
# goroutine 1
# bt
# up 10
