package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article PublishedArticle) error
	GetById(ctx context.Context, id int64) (Article, error)
	//SyncV2(ctx context.Context, article domain.Article) (int64, error)
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (G *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := G.db.WithContext(ctx).
		Where("id = ?", id).First(&art).Error
	return art, err
}

func (G *GORMArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	//return
	//panic("implement me")
	//先操作制作库，再操作线上库
	var (
		id  = article.Id
		err error
	)

	err = G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(tx)
		//publishArt := PublishedArticle(article)
		if id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		return txDAO.Upsert(ctx, PublishedArticle{Article: article})
	})
	return 0, err
}

func (G *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := G.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (G *GORMArticleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := G.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		}),
	}).Create(&art).Error
	return err
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
	Status   uint8 `gorm:"type:int(11);not null"`
}

//type PublishedArticle Article

type PublishedArticle struct {
	Article
}
