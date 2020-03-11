-- redo
-- ARGV: TIMEOUT

local erase = function(taskId)
    redis.call('HDEL', 'task-timestamp-hset', taskId);
    redis.call('HDEL', 'task-tracer-hset', taskId);
    redis.call('HDEL', 'task-result-hset', taskId);
end

local push = function(taskId)
    redis.call('HSET', 'task-stat-hset', taskId, 'PENDING');
    redis.call('LPUSH', 'task-pending-list', taskId);
end

-- traveral task in doing
local timeout = tonumber(ARGV[1])
local redo = {};
local now = redis.call('time')[1];
local tasks = redis.call('HGETALL', 'task-timestamp-hset');

for i = 1, #tasks, 2 do
    local duration = now - tasks[i + 1]
    if type(duration) == 'number' and duration > timeout then
        erase(tasks[i]);
        table.insert(redo, tasks[i]);
    end
end

-- insert to pending list
for i = 1, #redo do
    push(redo[i]);
end

return 1
