package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	logger2 "gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Data: float64(1),
				Msg:  "ok",
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish error"))
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			//模拟登陆态
			server.Use(func(c *gin.Context) {
				c.Set("claims", &ijwt.UserClaims{
					Uid: 123,
				})
			})
			h := NewArticleHandler(tt.mock(ctrl), &logger2.NopLogger{})
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			//这里加上了 可以直接往下走
			resp := httptest.NewRecorder()
			//t.Log(resp)
			server.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
