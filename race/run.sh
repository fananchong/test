#!/bin/bash

nohup ./race 2> >(./json.sh | tee -a 1.log) &

docker rm -f xxx
docker run --name='xxx' test_race