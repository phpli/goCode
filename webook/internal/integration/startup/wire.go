//go:build wireinject

package startup

import (
	//"gitee.com/geekbang/basic-go/webook/internal/events/article"
	//"gitee.com/geekbang/basic-go/webook/internal/job"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	//"gitee.com/geekbang/basic-go/webook/internal/service/sms/async"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
	//InitSaramaClient,
	//InitSyncProducer,
	InitLogger)

//var jobProviderSet = wire.NewSet(
//	service.NewCronJobService,
//	repository.NewPreemptJobRepository,
//	dao.NewGORMJobDAO)

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)

//var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO,
//	cache.NewInteractiveRedisCache,
//	repository.NewCachedInteractiveRepository,
//	service.NewInteractiveService,
//)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewCodeService,
		//service.NewUserService,
		service.NewArticleService,
		//InitWechatService,

		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WeChatHandler,
		web.NewArticleHandler,
		NewWechatHandlerConfig,
		InitWechatService,
		ijwt.NewRedisJWTHandler,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

//func InitAsyncSmsService(svc sms.Service) *async.Service {
//	wire.Build(thirdPartySet, repository.NewAsyncSMSRepository,
//		dao.NewGORMAsyncSmsDAO,
//		async.NewService,
//	)
//	return &async.Service{}
//}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		//userSvcProvider,
		service.NewArticleService,
		web.NewArticleHandler,
		//interactiveSvcSet,
		//repository.NewCachedArticleRepository,
		//cache.NewArticleRedisCache,
		//article.NewSaramaSyncProducer,
	)
	return &web.ArticleHandler{}
}

//func InitInteractiveService() service.InteractiveService {
//	wire.Build(thirdPartySet, interactiveSvcSet)
//	return service.NewInteractiveService(nil)
//}

//func InitJobScheduler() *job.Scheduler {
//	wire.Build(jobProviderSet, thirdPartySet, job.NewScheduler)
//	return &job.Scheduler{}
//}
