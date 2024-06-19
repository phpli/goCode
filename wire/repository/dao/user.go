package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func (u2 UserDao) Insert(ctx context.Context, u interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (u2 UserDao) Update(ctx context.Context, u interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (u2 UserDao) FindByEmail(ctx context.Context, email string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (u2 UserDao) FindByPhone(ctx context.Context, phone string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (u2 UserDao) FindById(ctx context.Context, id int64) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (u2 UserDao) UpdateById(ctx context.Context, entity interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
