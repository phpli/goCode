package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"gitee.com/geekbang/basic-go/webook/internal/web"
)

func InitWechatService() wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_ID ")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	// 692jdHsogrsYqxaUK9fgxw
	return wechat.NewService("692jdHsogrsYqxaUK9fgxw", "")
}

//func NewWechatHandlerConfig() web.WechatHandlerConfig {
//	return web.WechatHandlerConfig{
//		Secure: false,
//	}
//}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
