//go:build wireinject

package integration

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewCodeRepository,
		repository.NewCachedUserRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		//gin.Default,
	)
	return new(gin.Engine)
}
