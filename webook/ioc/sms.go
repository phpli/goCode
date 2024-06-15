package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/mermory"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)

func InitSMSService() sms.Service {
	return mermory.NewService()
}
