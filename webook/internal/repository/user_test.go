package repository

import (
	"context"
	"database/sql"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	cachemocks "gitee.com/geekbang/basic-go/webook/internal/repository/cache/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "12@qq.com",
						Valid:  true,
					},
					Password: "密码",
					Phone: sql.NullString{
						String: "1511112211",
						Valid:  true,
					},
					Nickname:    "li",
					Description: "1122111",
					Birthday:    now.UnixMilli(),
					Gender:      0,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:          123,
					Email:       "12@qq.com",
					Password:    "密码",
					Phone:       "1511112211",
					Nickname:    "li",
					Description: "1122111",
					Birthday:    now,
					Gender:      0,
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:          123,
				Email:       "12@qq.com",
				Password:    "密码",
				Phone:       "1511112211",
				Nickname:    "li",
				Description: "1122111",
				Birthday:    now,
				Gender:      0,
			},
			wantErr: nil,
		},
		{
			name: "缓存查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:          123,
					Email:       "12@qq.com",
					Password:    "密码",
					Phone:       "1511112211",
					Nickname:    "li",
					Description: "1122111",
					Birthday:    now,
					Gender:      0,
				}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:          123,
				Email:       "12@qq.com",
				Password:    "密码",
				Phone:       "1511112211",
				Nickname:    "li",
				Description: "1122111",
				Birthday:    now,
				Gender:      0,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, errors.New("mock error"))
				return d, c
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tt.mock(ctrl)
			repo := NewCachedUserRepository(ud, uc)
			u, err := repo.FindById(tt.ctx, tt.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUser, u)
			time.Sleep(time.Second) //测试gogo func
		})
	}
}
