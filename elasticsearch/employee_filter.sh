#!/bin/bash

set -x

# 筛选

curl -XGET 'http://localhost:9200/megacorp/employee/_search?pretty' -H "Content-Type: application/json" -d '{
    "query": {
        "bool": {
            "must": {
                "match": {
                    "last_name": "Smith"
                }
            },
            "filter": {
                "range": {
                    "age": { "gt": 30 }
                }
            }
        }
    }
}
'
