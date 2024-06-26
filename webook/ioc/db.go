package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	logger2 "gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(l logger2.LoggerV1) *gorm.DB {
	//dsn := viper.GetString("db.mysql.dsn")
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default", // yaml 文件里没有的话，会读这个默认值
	}
	//err := viper.UnmarshalKey("db.mysql", &cfg) db.mysql 不能带db.
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			//慢查询50ms，100ms. 一般 一次的磁盘io是10ms，给50ms 一般就是 一句给了5次
			SlowThreshold:             time.Millisecond * 50, //慢查询日志的阈值
			IgnoreRecordNotFoundError: true,                  //dev 环境设置为 false 比较好，正式环境设置为true
			LogLevel:                  glogger.Info,
			ParameterizedQueries:      true,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger2.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger2.Field{Key: "args", Value: args})
}
