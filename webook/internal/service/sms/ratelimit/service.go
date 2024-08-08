package ratelimit

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发限流")

type SmsService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, l ratelimit.Limiter) sms.Service {
	return &SmsService{
		svc:     svc,
		limiter: l,
	}
}

func (s *SmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		//限流器系统错误
		//可以限流：后续的接口很坑
		//可以不限：你的下游很强，业务可用性要求很高，尽量容错
		return fmt.Errorf("短信服务判断是否限流出问题", err)
	}
	if limited {
		return errLimited
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}
