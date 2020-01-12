#!/bin/sh

open http://localhost:9092/api/hello

docker run -it --name weesvc-gorilla --rm -p 9092:9092 \
 github.com/weesvc/weesvc-gorilla:0.0.1-SNAPSHOT \
 /bin/sh -c "/app/weesvc migrate; /app/weesvc serve"
