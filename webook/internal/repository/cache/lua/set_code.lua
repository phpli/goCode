--你的验证码在redis上的key
--phone_code:login:186xxxxx
local key = KEYS[1]
--验证次数，我们一个验证码，最多重复三次，这个记录了验证了几次,还可以验证几次
--phone_code:login:186xxxxx:cnt
local cntKey = key..":cnt"
--你的验证码 123456
local val = ARGV[1]
--过期时间
local ttl = tonumber(redis.call("ttl",key))

if ttl == -1 then
--key存在但是没有过期时间,手残删除了
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    return 0
    --完美
else
    --发送太频繁
    return -1
end


