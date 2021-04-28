#!/bin/bash

set -x

# 高亮

curl -XGET 'http://localhost:9200/megacorp/employee/_search?pretty' -H "Content-Type: application/json" -d '{
    "query" : {
        "match_phrase": {
            "about": "rock climbing"
        }
    },
    "highlight" : {
        "fields": {
            "about": {}
        }
    }
}
'
