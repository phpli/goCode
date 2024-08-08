package article

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	dao "gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	//同步存储
	SyncV1(ctx context.Context, article domain.Article) (int64, error)
	Sync(ctx context.Context, article domain.Article) (int64, error)
	//SyncV2(ctx context.Context, article domain.Article) (int64, error)
}
type CachedArticleRepository struct {
	dao dao.ArticleDAO

	//v1操作2个
	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO
	//db        *gorm.DB //事物 应该尽量在dao层
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   uint8(article.Status),
	})
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	//tx := c.db.WithContext(ctx).Begin()
	//if tx.Error != nil {
	//	return 0, tx.Error
	//}
	return c.dao.Sync(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	artn := c.toEntity(article)
	if id > 0 {
		err = c.authorDAO.UpdateById(ctx, artn)
	} else {
		id, err = c.authorDAO.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}