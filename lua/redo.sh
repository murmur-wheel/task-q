#!/bin/bash
TIMEOUT=1
function prepare() {
  redis-cli flushdb
  redis-cli --eval add-task.lua , task1 data1
  redis-cli --eval select-task.lua , tracer
}

prepare >/dev/null

seconds=$(($TIMEOUT + 1))
echo ">> sleep '$seconds's"
sleep $seconds

redis-cli --eval redo-expired-task.lua , $TIMEOUT
