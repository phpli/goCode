package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_sendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name    string
		before  func(t *testing.T) //准备数据
		after   func(t *testing.T) //验证数据
		reqBody string
		//redis
		//数据库
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15154107793").Result()
				defer cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
		"phone":"15154107793"

}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送频繁",
			before: func(t *testing.T) {
				//这个手机号，以及有一个验证码
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15154107793", "123456", time.Minute*9+time.Second*30).Result()
				defer cancel()
				assert.NoError(t, err)
				//assert.True(t, len(val) == 6)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15154107793").Result()
				defer cancel()
				assert.NoError(t, err)
				assert.True(t, "123456" == val)
				assert.Equal(t, "123456", val)
			},
			reqBody: `
{
		"phone":"15154107793"

}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送太频繁,稍后重试",
			},
		},
		{
			name: "数据格式有误",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			reqBody: `
{
		"phone":'15154107793'

}`,
			wantCode: http.StatusBadRequest,
			wantBody: web.Result{
				Msg: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			if resp.Code != 200 {
				return
			}
			assert.Equal(t, tc.wantCode, resp.Code)
			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
			//tc.after(t)
		})
	}
}
