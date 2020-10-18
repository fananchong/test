#!/bin/bash

set -e

go build

docker build -t virtual_memory:latest .
