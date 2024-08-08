package ratelimit

import "context"

type Limiter interface {
	//limited 有没有触发限流。key是限流对象
	//bool 是否限流
	//error 限流器是否有错误
	Limit(ctx context.Context, key string) (bool, error)
}
