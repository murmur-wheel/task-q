#!/bin/bash
# build image
# shellcheck disable=SC2046
docker run -p 6379:6379 -v $(pwd):/data/lua --name my-redis -d redis
