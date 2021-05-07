#!/bin/bash

time=$(date "+%Y-%m-%d %H:%M:%S")
while read line; do
    echo "${time}: ${line}"
done
