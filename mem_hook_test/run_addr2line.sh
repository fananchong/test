#!/bin/bash

set -ex

addr2line -e ./out $1 -f -a -p -C
