package web

import (
	"bytes"
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
		//fields fields
		//args   args
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12@qq.com",
					Password: "aa12@qqcom",
				}).Return(nil)
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "aa12@qqcom","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "邮箱错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "12qq.com",
				//	Password: "aa12@qqcom",
				//}).Return(nil)
				return usersvc
			},
			reqBody:  `{"email": "12qq.com","password": "aa12@qqcom","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "邮箱错误",
		},
		{
			name: "参数json不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "12qq.com",
				//	Password: "aa12@qqcom",
				//}).Return(http.StatusBadRequest)
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "a12@_com","ConfirmPassword": "aa12@qqcom";;;""11`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "密码格式有误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "12qq.com",
				//	Password: "aa12@qqcom",
				//}).Return(http.StatusBadRequest)
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "aa12@_com","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusOK,
			wantBody: "密码格式有误",
		},
		{
			name: "两次输入不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "aa12@qq1com","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusOK,
			wantBody: "两次输入不一致",
		}, {
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12@qq.com",
					Password: "aa12@qqcom",
				}).Return(service.ErrUserDuplicate)
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "aa12@qqcom","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12@qq.com",
					Password: "aa12@qqcom",
				}).Return(errors.New("系统错误"))
				return usersvc
			},
			reqBody:  `{"email": "12@qq.com","password": "aa12@qqcom","ConfirmPassword": "aa12@qqcom"}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			usersvc := tt.mock(ctrl)
			h := NewUserHandler(usersvc, nil)
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			//这里加上了 可以直接往下走
			resp := httptest.NewRecorder()
			//t.Log(resp)
			server.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
		Email: "test@test.com",
	}).Return(errors.New("test1 error"))
	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "test@test.com",
	})
	t.Log(err)
}

func TestUserHandler_loginSms(t *testing.T) {

}
