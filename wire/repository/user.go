package repository

import "gitee.com/geekbang/basic-go/wire/repository/dao"

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}
