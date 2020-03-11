-- add new task
-- ARGV: taskId, taskData
redis.call('RPUSH', 'task-pending-list', ARGV[1]);
redis.call('HSET', 'task-data-hset', ARGV[1], ARGV[2]);
redis.call('HSET', 'task-stat-hset', ARGV[1], 'PENDING');

return 1 -- redis.int