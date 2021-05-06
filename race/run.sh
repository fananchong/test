#!/bin/bash

./race 2> >(./date.sh | tee -a 1.log)

docker rm -f xxx
docker run --name='xxx' test_race