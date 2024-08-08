local key = KEYS[1]
--和用户的输入的code
local expectedCode = ARGV[1]
local cntKey = key..":cnt"
local code = redis.call("get",key)
--转成一个数字
local cnt = tonumber(redis.call("get",cntKey)) or 0
if cnt == nil or cnt <= 0 then
    -- 说明用户一直在输错
    return -1
elseif expectedCode == code then
    --正确,用完不能再用了
    redis.call("set",cntKey,-1)
    return 0
else
    --手斗输错
    redis.call("decr",cntKey)
    return -2
end


