package main

import (
	"gitee.com/geekbang/basic-go/webook/config"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/mermory"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server := initWebServer(redisClient)
	initUser(server, db, redisClient)
	server.Run(":8080")

	//练习部署用
	//server := gin.Default()
	//server.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello world k8s")
	//})
	//server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(server *gin.Engine, db *gorm.DB, rdb redis.Cmdable) {
	ud := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, userCache)
	us := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codePepo := repository.NewCodeRepository(codeCache)
	smsSvc := mermory.NewService()
	codeSvc := service.NewCodeService(codePepo, smsSvc)
	c := web.NewUserHandler(us, codeSvc)
	c.RegisterRoutes(server)
}

func initWebServer(redis *redis.Client) *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
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
	}))
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redis, time.Minute, 100).Build())
	//store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16, "tcp", "localhost:16379", "", []byte("fb0e22c79ac75679e9881e6ba183b354"),
	//	[]byte("988782dc147d58ff394f19a0d468d5b2"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("webook", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())
	return server
}

//func useSession(server *gin.Engine) {
//	login := &middleware.LoginMiddlewareBuilder{}
//	// 存储数据的，也就是你 userId 存哪里
//	// 直接存 cookie
//	store := cookie.NewStore([]byte("secret"))
//	// 基于内存的实现
//	//store := memstore.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("eF1`yQ9>yT1`tH1,sJ0.zD8;mZ9~nC6("))
//	//store, err := redis.NewStore(16, "tcp",
//	//	"localhost:6379", "",
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
//
//func useJWT(server *gin.Engine) {
//	login := middleware.LoginJWTMiddlewareBuilder{}
//	server.Use(login.CheckLogin())
//}
