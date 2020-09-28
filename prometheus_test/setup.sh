#!/bin/bash

CFG_DIR=$PWD/cfg
DATA_DIR=$PWD/data
echo $CFG_DIR

docker rm -f prometheus
docker rm -f node-exporter
docker rm -f grafana
docker rm -f alertmanager

docker run -d --name node-exporter --restart=always -p 9100:9100 -v "/:/host:ro,rslave" quay.io/prometheus/node-exporter --path.rootfs=/host
docker run -d --name grafana --restart=always --user root -p 3000:3000 -v $DATA_DIR/grafana:/var/lib/grafana grafana/grafana
docker run -d --name alertmanager --restart=always -p 9093:9093 -v $CFG_DIR:/etc/alertmanager prom/alertmanager

docker run -d --name prometheus --restart=always \
  --user root \
  -p 9090:9090 \
  -v $CFG_DIR:/etc/prometheus \
  -v $DATA_DIR/prometheus:/prometheus \
  --link node-exporter \
  --link alertmanager \
  prom/prometheus:latest

sleep 1s

docker ps

