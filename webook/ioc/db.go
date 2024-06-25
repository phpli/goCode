package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
