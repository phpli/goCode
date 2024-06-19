package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//输入
		ctx      context.Context
		email    string
		password string
		//输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "12@qq.com").Return(
					domain.User{
						Email:    "12@qq.com",
						Password: "$2a$10$wFt7cbZlJu95vcSurb5n8.BcfRk78ZlMJnJe.Jb2hxqSawQiRW02W",
						Phone:    "2222",
						Ctime:    now,
					}, nil)
				return repo
			},
			ctx:      context.Background(),
			email:    "12@qq.com",
			password: "aa12@qqcom",
			wantUser: domain.User{
				Email:    "12@qq.com",
				Password: "$2a$10$wFt7cbZlJu95vcSurb5n8.BcfRk78ZlMJnJe.Jb2hxqSawQiRW02W",
				Phone:    "2222",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "record not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "12@qq.com").Return(
					domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			ctx:      context.Background(),
			email:    "12@qq.com",
			password: "aa12@qqcom",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(
					domain.User{}, errors.New("mock db错误"))
				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "aa12@qqcom",
			wantUser: domain.User{},
			wantErr:  errors.New("mock db错误"),
		},
		{
			name: "用户名密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "12@qq.com").Return(
					domain.User{}, ErrInvalidUserOrPassword)
				return repo
			},
			ctx:      context.Background(),
			email:    "12@qq.com",
			password: "aaa12@qqcom",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

//func TestEncrypted(t *testing.T) {
//	res, err := bcrypt.GenerateFromPassword([]byte("aa12@qqcom"), bcrypt.DefaultCost)
//	if err == nil {
//		t.Log(string(res))
//	}
//}
