package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaSlideScript string

type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	//窗口大小
	interval time.Duration
	// interval内允许的阈值
	rate int
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	//key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return r.cmd.Eval(ctx, luaSlideScript, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
