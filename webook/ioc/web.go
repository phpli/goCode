package ioc

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middlewares/logger"
	logger2 "gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2WeChatHdl *web.OAuth2WeChatHandler,
	articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2WeChatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redis redis.Cmdable, l logger2.LoggerV1, jwtHdl ijwt.Handler) []gin.HandlerFunc {

	bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
		l.Debug("http 请求", logger2.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowRespBody(true)
	viper.OnConfigChange(func(e fsnotify.Event) { //监听配置文件的修改
		ok := viper.GetBool("web.logreq")
		bd.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		corsHDL(),
		bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/oauth2/wechat/authrul").
			IgnorePaths("/users/login_sms").Build(),
		//ratelimit.NewBuilder(redis, time.Minute, 100).Build(),
	}
}

func corsHDL() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost"},
		//是否允许带cookie
		AllowCredentials: true,
		//不写就是全部
		//AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			//
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
