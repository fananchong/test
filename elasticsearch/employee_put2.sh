#!/bin/bash

set -x

curl -XPUT 'http://localhost:9200/megacorp/employee/4?pretty' -H "Content-Type: application/json" -d '{
	"first_name2" : "John",
	"last_name2" : "Smith",
	"age2" : 25,
	"about2" : "I love to go rock climbing",
	"interests2": [ "sports", "music" ]
}
'

curl -XPUT 'http://localhost:9200/megacorp/employee/5?pretty' -H "Content-Type: application/json" -d '{
	"first_name2" : "Jane",
	"last_name2" : "Smith",
	"age2" : 32,
	"about2" : "I like to collect rock albums",
	"interests2": [ "music" ]
}
'

curl -XPUT 'http://localhost:9200/megacorp/employee/6?pretty' -H "Content-Type: application/json" -d '{
	"first_name2" : "Douglas",
	"last_name2" : "Fir",
	"age2" : 35,
	"about2" : "I like to build cabinets",
	"interests2": [ "forestry" ]
}
'

curl -XDELETE 'http://localhost:9200/megacorp/employee/4?pretty' -H "Content-Type: application/json"
curl -XDELETE 'http://localhost:9200/megacorp/employee/5?pretty' -H "Content-Type: application/json"
curl -XDELETE 'http://localhost:9200/megacorp/employee/6?pretty' -H "Content-Type: application/json"