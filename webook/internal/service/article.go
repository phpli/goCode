package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
}

func NewArticleService() ArticleService {
	return &articleService{}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	//TODO implement me
	return 1, nil
}
