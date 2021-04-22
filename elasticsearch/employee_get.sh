#!/bin/bash

set -x

curl -XGET 'http://localhost:9200/megacorp/employee/1?pretty' -H "Content-Type: application/json"

curl -XGET 'http://localhost:9200/megacorp/employee/2?pretty' -H "Content-Type: application/json"

curl -XGET 'http://localhost:9200/megacorp/employee/3?pretty' -H "Content-Type: application/json"
