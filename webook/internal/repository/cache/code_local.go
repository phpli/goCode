package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

type localCodeCache struct {
	cache      *lru.Cache[string, any]
	lock       sync.Mutex
	expiration time.Duration
	maps       sync.Map
}

func (l *localCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	//l.lock.Lock()
	//defer l.lock.Unlock()
	//key 非常多，maps占据很多内存
	key := l.key(biz, phone)
	lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//defer lock.(*sync.Mutex).Unlock()

	//换种写法
	defer func() {
		l.maps.Delete(key)
		lock.(*sync.Mutex).Unlock()
	}()

	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		l.cache.Add(key, codeItem{
			expire: now.Add(l.expiration),
			code:   code,
			cnt:    3,
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if itm.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSendTooMany
	}
	l.cache.Add(key, codeItem{
		expire: now.Add(l.expiration),
		code:   code,
		cnt:    3,
	})
	return nil
}

func (l *localCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	//l.lock.Lock()
	//defer l.lock.Unlock()

	key := l.key(biz, phone)
	lock, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	//defer lock.(*sync.Mutex).Unlock()

	//换种写法
	defer func() {
		l.maps.Delete(key)
		lock.(*sync.Mutex).Unlock()
	}()
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrKeyNotExist
	}
	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	return itm.code == inputCode, nil
}

func newLocalCodeCache(cache *lru.Cache[string, any], expiration time.Duration) *localCodeCache {
	return &localCodeCache{
		expiration: expiration,
		cache:      cache,
	}
}

func (l *localCodeCache) key(Biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", Biz, phone)
}

type codeItem struct {
	code string
	// 可验证次数
	cnt int
	// 过期时间
	expire time.Time
}
