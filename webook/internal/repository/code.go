package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (r *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return r.cache.Set(ctx, biz, phone, code)
}

func (r *CodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, inputCode)

}
