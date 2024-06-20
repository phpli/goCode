package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFalioverService struct {
	//你的服务商
	svcs []sms.Service
	idx  int32
	//连续超时个数
	cnt int32
	//阈值 超过多少个就要切换
	threshold int32
}

func NewTimeoutFalioverService(svcs []sms.Service) sms.Service {
	return &TimeoutFalioverService{
		svcs: svcs,
	}
}

func (t *TimeoutFalioverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//我成功往后挪一位
			atomic.StoreInt32(&t.cnt, 0)
		}
		//idx = newIdx
		idx = atomic.LoadInt32(&t.idx)
		svc := t.svcs[idx]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			atomic.AddInt32(&t.cnt, 1)
			return err
		case err == nil:
			atomic.AddInt32(&t.cnt, 0)
			return nil
		default:
			//不知道什么错误，
			//你可以考虑换一个
			//超时，可能是偶发的，我尽量试试
			//非超时，我直接下一个
			return err
		}
	}
	return nil
}
