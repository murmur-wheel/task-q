#!/bin/bash
# show data

echo ">> task-pending-list"
redis-cli lrange task-pending-list 0 -1
echo

echo ">> task-stat-hset"
redis-cli hgetall task-stat-hset
echo

echo ">> task-timestamp-hset"
redis-cli hgetall task-timestamp-hset
echo

echo ">> task-tracer-hset"
redis-cli hgetall task-tracer-hset
echo

echo ">> task-data-hset"
redis-cli hgetall task-data-hset
echo


echo ">> task-result-hset"
redis-cli hgetall task-result-hset
echo