#!/bin/bash

set -e

go build

docker build -t big_file:latest .
