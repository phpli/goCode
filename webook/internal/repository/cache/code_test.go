package cache

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
		//输入
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))

				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"1234879"}).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "1234879",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				//res.SetVal(int64(0))
				res.SetErr(errors.New("mock error"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"1234879"}).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "1234879",
			wantErr: errors.New("mock error"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"1234879"}).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "1234879",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-10))
				res.SetErr(errors.New("系统错误"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"1234879"}).Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "1234879",
			wantErr: errors.New("系统错误"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
