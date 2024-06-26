package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
)

func InitUserHandler(repo repository.UserRepository, logger logger.LoggerV1) service.UserService {

	return service.NewUserService(repo, logger)
}
