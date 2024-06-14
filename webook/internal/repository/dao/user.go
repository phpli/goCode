package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或者手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *UserDAO) Update(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	err := dao.db.WithContext(ctx).Save(&u).Error
	return err
}

type User struct {
	Id          int64          `gorm:"primaryKey,autoIncrement"`
	Email       sql.NullString `gorm:"type:varchar(255);unique"`
	Nickname    string         `gorm:"type:varchar(255)"`
	Phone       sql.NullString `gorm:"type:varchar(255);unique"`
	Password    string
	Ctime       int64
	Utime       int64
	Birthday    int64
	Gender      int
	Description string
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	//err := dao.db.WithContext(ctx).Where("email=?",email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	//err := dao.db.WithContext(ctx).Where("email=?",email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "id = ?", id).Error
	//err := dao.db.WithContext(ctx).Where("email=?",email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, entity User) error {

	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"utime":       time.Now().UnixMilli(),
			"nickname":    entity.Nickname,
			"birthday":    entity.Birthday,
			"description": entity.Description,
			"gender":      entity.Gender,
		}).Error
}
