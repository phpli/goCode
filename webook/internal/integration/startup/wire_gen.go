// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	loggerV1 := InitLogger()
	handler := jwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitMiddlewares(cmdable, loggerV1, handler)
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, loggerV1)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, cmdable, handler)
	wechatService := InitWechatService(loggerV1)
	wechatHandlerConfig := NewWechatHandlerConfig()
	oAuth2WeChatHandler := web.NewOAuth2WeChatHandler(wechatService, userService, wechatHandlerConfig, handler)
	articleDAO := dao.NewGORMArticleDAO(db)
	articleRepository := repository.NewArticleRepository(articleDAO)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WeChatHandler, articleHandler)
	return engine
}

func InitArticleHandler() *web.ArticleHandler {
	db := InitDB()
	articleDAO := dao.NewGORMArticleDAO(db)
	articleRepository := repository.NewArticleRepository(articleDAO)
	articleService := service.NewArticleService(articleRepository)
	loggerV1 := InitLogger()
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	return articleHandler
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,

	InitLogger)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewCachedUserRepository, service.NewUserService)
