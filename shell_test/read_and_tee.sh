#!/bin/bash

ls . | read
echo "cmd return value: "$?
ls . > ~/aa.log | read
echo "cmd return value: "$?
ls . | tee /dev/stderr |  read
echo "cmd return value: "$?
