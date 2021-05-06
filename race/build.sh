#!/bin/bash

go build -race -ldflags '-linkmode "external" -extldflags "-static"' .

docker build -t test_race:latest .