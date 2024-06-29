package integration

import (
	"bytes"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ArticleTestSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	//在测试执行前，初始化内容
	//s.server = startup.InitWebServer()

	//另一种方式,注册路由
	s.server = gin.Default()
	s.server.Use(func(c *gin.Context) {
		c.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	s.db = startup.InitDB()
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		//准备数据
		before func(t *testing.T)
		//验证数据
		after func(t *testing.T)
		//预期输入参数
		art Article
		//http
		wantCode int
		//预期输出
		wantRes Result[int64]
	}{
		{
			name: "新建帖子",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("title = ?", "我的标题").First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "ok",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求
			//执行请求
			//验证结果
			tc.before(t)
			defer tc.after(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			if resp.Code != 200 {
				return
			}
			assert.Equal(t, tc.wantCode, resp.Code)
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}

// 每次test执行
func (s *ArticleTestSuite) TearDownSuite() {
	s.db.Exec("TRUNCATE TABLE articles")
}
