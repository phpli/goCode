package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	//db := initDB()
	//redisClient := initRedis()
	//
	//server := initWebServer(redisClient)
	//
	//initUser(server, db, redisClient)
	initViperV1()
	//initViperRemote()
	//initViperByConfig()
	server := InitWebServer()
	server.Run(":8080")

	//练习部署用
	//server := gin.Default()
	//server.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello world k8s")
	//})
	//server.Run(":8080")
}

//func initRedis() redis.Cmdable {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	return redisClient
//}
//
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
//
//func initUser(server *gin.Engine, db *gorm.DB, rdb redis.Cmdable) {
//	ud := dao.NewUserDAO(db)
//	userCache := cache.NewUserCache(rdb)
//	repo := repository.NewUserRepository(ud, userCache)
//	us := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codePepo := repository.NewCodeRepository(codeCache)
//	smsSvc := mermory.NewService()
//	codeSvc := service.NewCodeService(codePepo, smsSvc)
//	c := web.NewUserHandler(us, codeSvc)
//	c.RegisterRoutes(server)
//}

//func initWebServer(redis redis.Cmdable) *gin.Engine {
//	server := gin.Default()
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins: []string{"http://localhost"},
//		//是否允许带cookie
//		AllowCredentials: true,
//		//不写就是全部
//		//AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
//		AllowHeaders:  []string{"Content-Type", "Authorization"},
//		ExposeHeaders: []string{"x-jwt-token"},
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				return true
//			}
//			//
//			return strings.Contains(origin, "your_company.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr: config.Config.Redis.Addr,
//	//})
//	server.Use(ratelimit.NewBuilder(redis, time.Minute, 100).Build())
//	//store := cookie.NewStore([]byte("secret"))
//	//store, err := redis.NewStore(16, "tcp", "localhost:16379", "", []byte("fb0e22c79ac75679e9881e6ba183b354"),
//	//	[]byte("988782dc147d58ff394f19a0d468d5b2"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	//server.Use(sessions.Sessions("webook", store))
//	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/users/signup").
//		IgnorePaths("/users/login").
//		IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		Build())
//	return server
//}

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

func initViper() {
	viper.SetConfigName("dev")      //文件名称
	viper.SetConfigType("yaml")     //文件后缀名
	viper.AddConfigPath("./config") //当前工作目录下的config子目录
	//viper.AddConfigPath("./tmp/config") //当前工作目录下的config子目录
	//viper.AddConfigPath("./etc/webook") //当前工作目录下的config子目录
	err := viper.ReadInConfig() //读取配置到内存
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//otherViper := viper.New()
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigName("myjson")
	//otherViper.SetConfigType("json")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViperV1() {
	viper.SetConfigFile("config/dev.yaml")
	viper.WatchConfig()
	//只告诉你文件变了
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name, e.Op)

	})
	err := viper.ReadInConfig() //读取配置到内存
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// 可读网络，可读预设 配置
func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:16379"
`

	err := viper.ReadConfig(bytes.NewBuffer([]byte(cfg)))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func initViperByConfig() {
	//设置页面 项目参数一列 --config=config/dev.yaml
	cfile := pflag.String("config", "config/config.yaml", "config file path")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig() //读取配置到内存
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
		//log.Println("watch", viper.GetString("test.key")) 有变化的远程配置，需要重新get一下
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
