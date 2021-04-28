#!/bin/bash

set -x

# 相关性搜索

curl -XGET 'http://localhost:9200/megacorp/employee/_search?pretty' -H "Content-Type: application/json" -d '{
    "query" : {
        "match": {
            "last_name": "Smith"
        }
    }
}
'
