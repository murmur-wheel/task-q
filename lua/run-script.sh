#!/bin/bash
if [ $# -eq 1 ] && [ -f $1 ]; then
  docker exec -it my-redis bash -c "cd /data/lua/lua; bash $1"
else
  echo "please specify a executable script"
fi
