#!/bin/bash

docker run -d --name=grafana -p 3000:3000 grafana/grafana

docker start grafana
