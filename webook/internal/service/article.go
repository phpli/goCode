package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/article"
	logger2 "gitee.com/geekbang/basic-go/webook/pkg/logger"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	//v  v 和v1是互斥的
	repo article.ArticleRepository

	//v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger2.LoggerV1
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository, l logger2.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}
func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}
	return a.repo.Create(ctx, article)
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	//if article.Id > 0 {
	//	err := a.repo.Update(ctx, article)
	//	return article.Id, err
	//}
	//return a.repo.Create(ctx, article)
	return a.repo.Sync(ctx, article)
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		err = a.author.Update(ctx, article)
	} else {
		id, err = a.author.Create(ctx, article)
		if err != nil {
			return 0, errors.New("mock db error")
		}
	}
	if err != nil {
		return 0, err
	}
	article.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, article)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败", logger2.String("article_id", id), logger2.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，重试彻底失败", logger2.String("article_id", id), logger2.Error(err))
	}
	return id, err
}
