#!/bin/bash

addr2line -e ./out $1 -f -a -p -C
