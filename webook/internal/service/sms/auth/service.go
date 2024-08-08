package auth

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

//send 一个是 []*string

func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	//是否是权限校验
	var tc Claims
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return s.svc.Send(ctx, biz, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
