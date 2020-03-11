#!/bin/bash
# push result
echo ">> flushdb"
redis-cli flushdb
echo

echo ">> prepare data"
redis-cli --eval add-task.lua , task1 data1
echo

echo ">> select a task"
tracerId="tracer1"
task=$(redis-cli --eval select-task.lua , $tracerId)
echo "receive task: $task"
echo

echo ">> push result"
redis-cli --eval push-task-result.lua , "$task" $tracerId "trace result"
