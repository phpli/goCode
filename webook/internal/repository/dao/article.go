package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (G *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := G.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// 建表语句
type Article struct {
	Id       int64  `gorm:"primary_key;auto_increment"`
	Title    string `gorm:"type:varchar(255);not null"`
	Content  string `gorm:"type:text;not null"`
	AuthorId int64  `gorm:"type:bigint(20);not null;index"`
	Ctime    int64
	Utime    int64
}
