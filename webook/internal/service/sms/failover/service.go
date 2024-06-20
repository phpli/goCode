package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"sync/atomic"
)

type FailoverService struct {
	svcs []sms.Service
	idx  uint64
}

func (f *FailoverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err != nil {
			return nil
		}
	}
	return errors.New("发送失败")
}

func (f *FailoverService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < (idx + length); i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return err
		}
	}
	return errors.New("发送失败")
}

func NewFailoverService(svcs []sms.Service) sms.Service {
	return &FailoverService{
		svcs: svcs,
	}
}
