//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		//本地
		DSN: "root:root@tcp(webook-mysql:13317)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:16389",
	},
}
