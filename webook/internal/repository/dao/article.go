package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
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

func (G *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	//依赖gorm的忽略0值的特性
	//err := G.db.WithContext(ctx).Updates(&art).Error
	res := G.db.WithContext(ctx).Model(&art).Where("id=? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		//补充日志 或者 监控+1
		return fmt.Errorf("更新失败，可能是创作者非法id=%d, author_id = %d", art.Id, art.AuthorId)
	}
	return res.Error
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
