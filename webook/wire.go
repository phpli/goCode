//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger,

		dao.NewUserDAO,
		dao.NewGORMArticleDAO,
		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewCodeRepository,
		repository.NewCachedUserRepository,
		repository.NewArticleRepository,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		ijwt.NewRedisJWTHandler,
		ioc.InitWechatService,
		ioc.InitSMSService,
		web.NewUserHandler,
		web.NewOAuth2WeChatHandler,
		web.NewArticleHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.NewWechatHandlerConfig,
		//gin.Default,
	)
	return new(gin.Engine)
}
