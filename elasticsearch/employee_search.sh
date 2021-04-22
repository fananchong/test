#!/bin/bash

set -x

curl -XGET 'http://localhost:9200/megacorp/employee/_search?pretty' -H "Content-Type: application/json"

curl -XGET 'http://localhost:9200/megacorp/employee/_search?q=last_name:Smith&pretty' -H "Content-Type: application/json"
