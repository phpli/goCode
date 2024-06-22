package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2WeChatHdl *web.OAuth2WeChatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2WeChatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redis redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHDL(),
		middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
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
		ExposeHeaders: []string{"x-jwt-token"},
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
