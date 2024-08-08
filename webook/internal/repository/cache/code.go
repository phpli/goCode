package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("send message too many")
	ErrCodeVerifyTooManyTimes = errors.New("code verify too many times")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
	//key(Biz string, phone string) string
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()

	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		//正常来说这里频繁出现，需要预警
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, nil
	}
}

func (c *RedisCodeCache) key(Biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", Biz, phone)
}
