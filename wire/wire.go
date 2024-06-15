//go:build wireinject

// 让wire 注入这里的代码
package wire

import (
	"gitee.com/geekbang/basic-go/wire/repository"
	"gitee.com/geekbang/basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	//这个方法传入各个组件的初始化
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
