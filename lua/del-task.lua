--
-- Created by IntelliJ IDEA.
-- User: murmur
-- Date: 3/11/20
-- Time: 5:17 AM
-- To change this template use File | Settings | File Templates.
--
local taskId = ARGV[1]
redis.call('HDEL', 'task-stat-hset', taskId)

return 1 -- int64