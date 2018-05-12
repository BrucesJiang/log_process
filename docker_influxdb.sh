#!/bin/bash

cd `dirname ${BASH_SOURCE-$0}`
#docker run -d -p 8083:8083 -p 8086:8086 --expose 8089 --expose 8099 --name influxdb tutum/influxdb

CPATH=`pwd`

echo $CPATH
docker run -itd -p 8083:8083 -p 8086:8086 -e ADMIN_USER="root" -e INFLUXDB_INIT_PWD="root" \
    -v influxdb:$CPATH/data --name influxdb tutum/influxdb:latest

#docker run --rm \
#      -e INFLUXDB_DB=log_process -e INFLUXDB_ADMIN_ENABLED=true \
#      -e INFLUXDB_ADMIN_USER=root -e INFLUXDB_ADMIN_PASSWORD=root \
#      -e INFLUXDB_USER=bruce -e INFLUXDB_USER_PASSWORD=bruce \
#      -v $PWD:$CPATH/data/influxdb \
#      influxdb ./init-influxdb.sh
echo "Run Backround"

#docker run -itd -p 8083:8083 -p 8086:8086 --name influxd influxdb:latest
docker start influxdb

echo "Start DB"
