#!/bin/bash

set -ex

openssl genrsa -out rsa.pem 1024
openssl rsa -in rsa.pem -pubout -out rsa.pub.pem
