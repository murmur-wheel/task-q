-- tracer acquire a new task
-- ARGV: tracerId
-- if succeed, true taskId, otherwise return ""
local tracerId = ARGV[1]

-- select trask into doing
local select = function(taskId, tracerId)
    -- update stat
    redis.call('HSET', 'task-stat-hset', taskId, 'DOING')
    -- tracer id
    redis.call('HSET', 'task-tracer-hset', taskId, tracerId);
    -- modify time
    local now = redis.call('time') -- UNIX timestamp
    redis.call('HSET', 'task-timestamp-hset', taskId, now[1]);
end

-- loop, until a task avail
while true
do
    local popped = redis.call('LPOP', 'task-pending-list')
    if type(popped) == 'boolean' then
        return "" -- there are not avail task anymore
    end

    if type(popped) == 'string' then
        local stat = redis.call('HGET', 'task-stat-hset', popped)
        if type(stat) == 'string' and stat == 'PENDING' then
            -- popped is the target
            -- select as a doing task
            select(popped, tracerId)
            return popped
        end
    end
end

