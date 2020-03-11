-- push result
-- KEYS:
-- ARGV: taskId tracerId result
-- if push success, return (integer)1, otherwise return (interger)0.
local stat = redis.call('HGET', 'task-stat-hset', ARGV[1]);

-- check stat
if type(stat) == 'string' and stat == 'DOING' then
    -- get current tracer
    local tracerId = redis.call('HGET', 'task-tracer-hset', ARGV[1]);

    -- verify tracer
    if type(tracerId) == 'string' and tracerId == ARGV[2] then
        -- modify timestamp
        local now = redis.call('time')[1]
        redis.call('HSET', 'task-timestamp-hset', ARGV[1], now);

        -- update result
        redis.call('HSET', 'task-result-hset', ARGV[1], ARGV[3]);
        return 1
    end
end

return 0
