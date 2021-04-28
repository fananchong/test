#!/bin/bash

set -x

# 聚合


curl -X PUT "localhost:9200/megacorp/_mapping?pretty" -H 'Content-Type: application/json' -d'
{
  "properties": {
    "interests": { 
      "type":     "text",
      "fielddata": true
    }
  }
}
'

curl -XGET 'http://localhost:9200/megacorp/employee/_search?pretty' -H "Content-Type: application/json" -d '{
    "aggs" : {
        "all_interests": {
            "terms": { "field": "interests" }
        }
    }
}
'

