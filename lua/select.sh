#!/bin/bash
# shellcheck disable=SC2164
echo ">> flushdb"
redis-cli flushdb

echo ">> prepare data"
redis-cli --eval add-task.lua , task1 data1
echo

echo ">> select a task"
redis-cli --eval select-task.lua , "tracerId1"
echo

echo ">> lrange task-pending-list"
redis-cli lrange task-pending-list 0 -1
echo

echo ">> hgetall task-timestamp-hset"
redis-cli hgetall task-timestamp-hset
echo

echo ">> hgetall task-stat-hset"
redis-cli hgetall task-stat-hset
